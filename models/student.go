package models

type Student struct {
	Id    string  `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Major string  `json:"major" binding:"required"`
	GPA   float64 `json:"gpa" binding:"required,min=0,max=4.00"`
}
