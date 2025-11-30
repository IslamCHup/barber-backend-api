package models

import "gorm.io/gorm"

type Barber struct {
	gorm.Model
	FullName       string  `json:"full_name" gorm:"not null"`
	WorkHoursStart int     `json:"work_hours_start" gorm:"not null;default:9"`
	WorkHoursEnd   int     `json:"work_hours_end" gorm:"not null;default:17"`
	AvgRating      float64 `json:"avg_rating" gorm:"not null;default:0"`
}

// type BarberUpdateReqDTO struct {
// 	WorkHoursStart int     `json:"work_hours_start" gorm:"not null;default:9"`
// 	WorkHoursEnd   int     `json:"work_hours_end" gorm:"not null;default:17"`
// }

type BarbersCreateReqDTO struct {
	FullName  string `json:"full_name" gorm:"not null"`
}

type BarberResDTO struct {
	FullName  string  `json:"full_name"`
	AvgRating float64 `json:"avg_rating"`
}
