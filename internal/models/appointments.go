package models

import (
	"gorm.io/gorm"
)

type Appointments struct {
	gorm.Model
	BarberID uint   `json:"barber_id" gorm:"not null"`
	Barber   Barber `json:"barber" gorm:"foreignKey:BarberID; not null"`
	ClientID uint   `json:"client_id" gorm:"not null"`
	Client   Client `json:"client" gorm:"foreignKey:ClientID; not null"`
	Time     string `json:"time" gorm:"not null"`
	Rating   *int   `json:"rating" gorm:"default:0"`
}

type AppointmentsCreateDTO struct {
	BarberID uint   `json:"barber_id" gorm:"not null"`
	ClientID uint   `json:"client_id" gorm:"not null"`
	Time     string `json:"time" gorm:"not null"`
}

type AppointmentsUpdateReqDTO struct {
	BarberID uint `json:"barber_id" gorm:"not null"`
	ClientID uint   `json:"client_id" gorm:"not null"`
	Rating   int  `json:"rating" gorm:"default:0"`
}
