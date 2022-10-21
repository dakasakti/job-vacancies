package entities

import "gorm.io/gorm"

type TimeBackend struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	BackEnds   []BackEnd `json:"back_ends"`
}

type BackEnd struct {
	gorm.Model    `json:"-"`
	ID            uint   `json:"id"`
	CompanyName   string `json:"company_name"`
	JobPosition   string `json:"job_position"`
	WorkType      string `json:"work_type"`
	TechStack     string `json:"tech_stack"`
	LinkToJob     string `json:"link"`
	Notes         string `json:"notes,omitempty"`
	Deadline      string `json:"deadline,omitempty"`
	TimeBackendID uint   `json:"time_id"`
}
