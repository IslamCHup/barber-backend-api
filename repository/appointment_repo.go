package repository

import (
	"barber-backend-api/internal/models"

	"gorm.io/gorm"
)

type AppointmentsRepository interface {
	GetAllAppointments() ([]models.Appointments, error)
	CreateAppointment(req *models.Appointments) error
	Update(id uint, req models.AppointmentsUpdateReqDTO) error
	GetLastAppointments(id uint) (*models.Appointments, error)
	GetAllAppointmentsByBarberID(id uint) ([]models.Appointments, error)
	GetByID(id uint) (*models.Appointments, error)
	Delete(id uint) error
}

type appointmentsRepository struct {
	db *gorm.DB
}

func NewAppointmentsRepository(db *gorm.DB) AppointmentsRepository {
	return &appointmentsRepository{db: db}
}

func (r *appointmentsRepository) GetAllAppointments() ([]models.Appointments, error) {
	var appointments []models.Appointments

	if err := r.db.Model(&models.Appointments{}).Preload("Barber").Preload("Client").Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}
func (r *appointmentsRepository) GetByID(id uint) (*models.Appointments, error) {
	var appnmts models.Appointments

	if err := r.db.Model(&models.Appointments{}).Preload("Barber").Preload("Client").First(&appnmts, id).Error; err != nil {
		return nil, err
	}

	return &appnmts, nil
}
func (r *appointmentsRepository) GetAllAppointmentsByBarberID(id uint) ([]models.Appointments, error) {
	var appnmtsBarber []models.Appointments
	if err := r.db.Model(&models.Appointments{}).Preload("Barber").Preload("Client").Where("barber_id = ?", id).Find(&appnmtsBarber).Error; err != nil {
		return nil, err
	}
	return appnmtsBarber, nil
}

func (r *appointmentsRepository) CreateAppointment(req *models.Appointments) error {
	if req == nil {
		return nil
	}
	return r.db.Create(req).Error
}
func (r *appointmentsRepository) Update(id uint, req models.AppointmentsUpdateReqDTO) error {
	return r.db.Model(&models.Appointments{}).Where("id = ?", id).Update("rating", req.Rating).Error
}

func (r *appointmentsRepository) GetLastAppointments(id uint) (*models.Appointments, error) {
	var lastAppointments models.Appointments
	if err := r.db.Model(&models.Appointments{}).Where("barber_id = ?", id).Last(&lastAppointments).Error; err != nil {
		return nil, err
	}
	return &lastAppointments, nil
}

func (r *appointmentsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Appointments{}, id).Error
}
