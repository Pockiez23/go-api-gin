```markdown
# 🎓 Student REST API (Go + Gin + SQLite)

![Go Version]
![Gin Framework]
![SQLite]

โปรเจกต์ RESTful API สำหรับจัดการข้อมูลนักเรียน (CRUD) พัฒนาด้วยภาษา **Go** โดยใช้ **Gin Framework** และฐานข้อมูล **SQLite** ออกแบบโครงสร้างแบบ Layered Architecture (Handler -> Service -> Repository)

---

## 🚀 ฟีเจอร์หลัก (Features)

* ✅ **Create:** เพิ่มข้อมูลนักเรียน (รองรับทั้งแบบ **คนเดียว** และ **หลายคนพร้อมกัน**)
* ✅ **Read:** ดึงข้อมูลนักเรียนทั้งหมด หรือค้นหาตาม ID
* ✅ **Update:** แก้ไขข้อมูลนักเรียน (ตรวจสอบ Validation)
* ✅ **Delete:** ลบข้อมูลนักเรียน
* 🛡️ **Error Handling:** จัดการ Error แบบมาตรฐาน ซ่อน Error จาก Database เพื่อความปลอดภัย

---

## 🛠️ โครงสร้างโปรเจกต์ (Project Structure)

```text
go-api-gin/
├── config/           # การตั้งค่าและการเชื่อมต่อ Database
├── handlers/         # รับ Request และส่ง Response (Controller)
├── models/           # โครงสร้างข้อมูล (Structs)
├── repositories/     # จัดการคำสั่ง SQL กับ Database โดยตรง
├── services/         # Business Logic ของระบบ
├── main.go           # จุดเริ่มต้นของโปรแกรม (Entry Point)
├── students.db       # ไฟล์ฐานข้อมูล SQLite
└── go.mod            # Go Modules dependency file

```

---

## ⚙️ การติดตั้งและเริ่มใช้งาน (Installation)

1. **Clone โปรเจกต์**
```bash
git clone [https://github.com/yourusername/student-api.git](https://github.com/yourusername/student-api.git)
cd student-api

```


2. **ติดตั้ง Library ที่จำเป็น**
```bash
go mod tidy

```


3. **รันเซิร์ฟเวอร์**
```bash
go run main.go

```


4. **เมื่อรันสำเร็จ** Terminal จะแสดงผลว่า:
```text
[GIN-debug] Listening and serving HTTP on :8080

```



---

## 📡 คู่มือการใช้งาน API (API Documentation)

### 1. ดึงข้อมูลนักเรียนทั้งหมด (Get All Students)

* **URL:** `/students`
* **Method:** `GET`

### 2. ดึงข้อมูลนักเรียนตาม ID (Get Student by ID)

* **URL:** `/students/:id`
* **Method:** `GET`
* **Example URL:** `http://localhost:8080/students/66090001`

---

### 3. เพิ่มข้อมูลนักเรียน (Create Student)

รองรับ 2 รูปแบบใน URL เดียวกัน!

* **URL:** `/students`
* **Method:** `POST`

#### แบบที่ A: เพิ่มทีละคน (Single Object)

```json
{
    "id": "66090001",
    "name": "Somchai Jaidee",
    "major": "CS",
    "gpa": 3.50
}

```

#### แบบที่ B: เพิ่มทีละหลายคน (Batch Array) 🌟

```json
[
    {
        "id": "66090002",
        "name": "Alice Wonderland",
        "major": "SE",
        "gpa": 3.80
    },
    {
        "id": "66090003",
        "name": "Bob Builder",
        "major": "IT",
        "gpa": 2.50
    }
]

```

---

### 4. แก้ไขข้อมูลนักเรียน (Update Student)

* **URL:** `/students/:id`
* **Method:** `PUT`
* **Validation:** GPA ต้องอยู่ระหว่าง 0.00 - 4.00

**Request Body:**

```json
{
    "name": "Somchai Updated",
    "major": "CS",
    "gpa": 3.90
}

```

---

### 5. ลบข้อมูลนักเรียน (Delete Student)

* **URL:** `/students/:id`
* **Method:** `DELETE`
* **Response:** `204 No Content` (ถ้าสำเร็จ) หรือ `404 Not Found` (ถ้าไม่เจอ)

---

## 🧪 วิธีทดสอบด้วย Postman (Testing)

1. เปิดโปรแกรม **Postman**
2. สร้าง Request ใหม่
3. เลือก Method (GET, POST, PUT, DELETE) ให้ตรงกับที่ต้องการ
4. ใส่ URL `http://localhost:8080/students`
5. หากเป็น **POST** หรือ **PUT**:
* ไปที่ tab **Body**
* เลือก **raw**
* เลือกรูปแบบเป็น **JSON**
* วาง JSON ตัวอย่างจากด้านบนลงไป


6. กดปุ่ม **Send** 🚀

---

## ⚠️ ปัญหาที่พบบ่อย (Troubleshooting)

* **Error: `bind: address already in use**`
* แปลว่า Port 8080 ถูกใช้อยู่ ให้ปิด Terminal เก่า หรือใช้คำสั่ง Kill process


* **Error: `database is locked**`
* โปรดตรวจสอบว่าคุณเปิดไฟล์ `students.db` ค้างไว้ในโปรแกรมอื่น (เช่น DB Browser) หรือไม่ ถ้าเปิดอยู่ให้ปิดโปรแกรมนั้นก่อน



---

Deverloped by ITTHICHET HONGWORAPAT 6509650187
