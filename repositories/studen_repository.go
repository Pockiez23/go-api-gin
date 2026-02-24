package repositories

import (
	"database/sql"
	"fmt"

	"example.com/student-api/models"
)

type StudentRepository struct {
	DB *sql.DB
}

func (r *StudentRepository) GetAll() ([]models.Student, error) {
	rows, err := r.DB.Query("SELECT id, name, major, gpa FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		rows.Scan(&s.Id, &s.Name, &s.Major, &s.GPA)
		students = append(students, s)
	}
	return students, nil
}

func (r *StudentRepository) GetByID(id string) (*models.Student, error) {
	row := r.DB.QueryRow(
		"SELECT id, name, major, gpa FROM students WHERE id = ?",
		id,
	)

	var s models.Student
	err := row.Scan(&s.Id, &s.Name, &s.Major, &s.GPA)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) Create(s models.Student) error {
	_, err := r.DB.Exec(
		"INSERT INTO students (id, name, major, gpa) VALUES (?, ?, ?, ?)",
		s.Id, s.Name, s.Major, s.GPA,
	)
	return err
}

func (r *StudentRepository) Update(id string, s models.Student) error {
	// คำสั่ง SQL เพื่อแก้ไขข้อมูล
	query := "UPDATE students SET name = ?, major = ?, gpa = ? WHERE id = ?"

	// Execute คำสั่ง
	result, err := r.DB.Exec(query, s.Name, s.Major, s.GPA, id)
	if err != nil {
		return err
	}

	// เช็คว่ามีแถวที่ถูกแก้ไขจริงไหม
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// ถ้า rowsAffected เป็น 0 แปลว่าหา ID ไม่เจอ -> ส่ง error กลับไป
	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *StudentRepository) Delete(id string) error {
	// สั่งลบข้อมูลตาม ID
	query := "DELETE FROM students WHERE id = ?"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// เช็คว่ามีแถวถูกลบจริงไหม
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// ถ้าไม่มีแถวไหนหายไปเลย แปลว่าหา ID ไม่เจอ
	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}
