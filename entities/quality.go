package entities

import "gorm.io/gorm"

type TimeQuality struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Qualitys   []Quality `json:"qualitys"`
}

type Quality struct {
	gorm.Model    `json:"-"`
	ID            uint   `json:"id"`
	CompanyName   string `json:"company_name"`
	JobPosition   string `json:"job_position"`
	WorkType      string `json:"work_type"`
	TechStack     string `json:"tech_stack"`
	LinkToJob     string `json:"link"`
	Notes         string `json:"notes,omitempty"`
	Deadline      string `json:"deadline,omitempty"`
	TimeQualityID uint   `json:"time_id"`
}
