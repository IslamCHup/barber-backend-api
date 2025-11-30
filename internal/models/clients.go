package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	FullName   string `json:"full_name" gorm:"not null"`
}

type ClientUpdateReqDTO struct {
	FullName   *string `json:"full_name"`
}

type ClientCreateReqDTO struct {
	FullName   string  `json:"full_name" gorm:"not null"`
}

type ClientRespDTO struct {
	FullName   string `json:"full_name"`
}
