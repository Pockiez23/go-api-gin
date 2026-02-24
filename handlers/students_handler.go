package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/student-api/models"
	"example.com/student-api/services"
)

type StudentHandler struct {
	Service *services.StudentService
}

// รับค่าทั้งหมดจาก Service แล้วส่งกลับเป็น JSON Response
func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.Service.GetStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

// รับค่า ID จาก URL Path และส่งต่อให้ Service เพื่อดึงข้อมูลนักเรียนตาม ID นั้นๆ
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

// รับข้อมูลนักเรียนจาก Request Body แล้วส่งต่อให้ Service เพื่อบันทึกลง Database
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	// อ่านข้อมูลดิบ (Bytes)
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// เตรียมตัวแปร Slice รับข้อมูล
	var students []models.Student

	// ตรวจสอบว่าเป็น Array หรือ Object
	trimmedBody := bytes.TrimSpace(bodyBytes)
	if len(trimmedBody) > 0 && trimmedBody[0] == '[' {
		// กรณี Array
		if err := json.Unmarshal(bodyBytes, &students); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format (Expected Array)"})
			return
		}
	} else {
		// กรณี Single Object
		var singleStudent models.Student
		if err := json.Unmarshal(bodyBytes, &singleStudent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format (Expected Object)"})
			return
		}
		students = append(students, singleStudent)
	}

	// วนลูปบันทึก + Validation + Error Hiding
	var successCount int
	var errorList []string

	for i, student := range students {
		// Validation แบบ Manual (เพราะ Unmarshal ไม่เช็ค binding tag ให้)
		if student.Name == "" || student.Major == "" || student.GPA < 0 || student.GPA > 4.00 {
			msg := fmt.Sprintf("Row %d (ID %v): ข้อมูลไม่ถูกต้อง (กรุณาเช็ค Name, Major หรือ GPA)", i+1, student.Id)
			errorList = append(errorList, msg)
			continue // ข้ามคนนี้ไปเลย
		}

		// บันทึกลง DB
		if err := h.Service.CreateStudent(student); err != nil {
			// Log Error จริงไว้ดูเองใน Terminal
			fmt.Printf("🚨 [Create Error!!!] ID %v: %v\n", student.Id, err)

			// ส่ง Error สุภาพๆ ให้ User (ซ่อน DB Error)
			// เช็คว่าเป็น Error เรื่อง ID ซ้ำหรือไม่ (Duplicate Key)
			errMsg := "Failed to save student data (Internal Error)"

			// หมายเหตุ: ตรงนี้ถ้าอยากให้ละเอียดขึ้น อาจจะเช็ค string err ว่ามีคำว่า "UNIQUE constraint" ไหม
			// แต่เพื่อความปลอดภัยสูงสุด ใช้ข้อความกลางๆ ไว้ก่อนครับ
			errorList = append(errorList, fmt.Sprintf("ID %v: %s", student.Id, errMsg))
		} else {
			successCount++
		}
	}

	// สรุปผลลัพธ์
	response := gin.H{
		"message":       "Batch processing completed",
		"received":      len(students),
		"success_count": successCount,
	}

	if len(errorList) > 0 {
		response["errors"] = errorList
		c.JSON(http.StatusMultiStatus, response)
	} else {
		c.JSON(http.StatusCreated, response) // 201 Created
	}
}

// PUT: UpdateStudent (Validation + Error Hiding)
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var student models.Student
	// Validation ผ่าน ShouldBindJSON (ใช้ Struct Tag binding:"required" ได้เลย)
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Input Data",
			"details": err.Error(), // (Validation Error ไม่อันตราย)
		})
		return
	}

	// กำหนด ID ให้ตรงกับ URL
	student.Id = id

	// เรียก Service
	err := h.Service.UpdateStudent(id, student)
	if err != nil {
		// แยกแยะ Error
		if err.Error() == "student not found" {
			// กรณี: หาไม่เจอ (User Fault) -> บอกตรงๆ ได้
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		} else {
			// กรณี: DB พัง / SQL ผิด (Server Fault) -> 🚨 Log + 🛡️ Hide
			fmt.Printf("🚨 [Update Error!!!] ID %s: %v\n", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error. Please try again later.",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"data":    student,
	})
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	// ดึง ID จาก URL
	id := c.Param("id")

	// สั่งลบผ่าน Service
	err := h.Service.DeleteStudent(id)
	if err != nil {
		// ถ้า Error เพราะหาไม่เจอ ให้ส่ง 404
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		} else {
			// Error อื่นๆ (เช่น DB พัง) ส่ง 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// ลบสำเร็จ! ส่ง 204 No Content (ไม่ต้องมี Body)
	c.Status(http.StatusNoContent)
}
