package repository

import (
	"barber-backend-api/internal/models"

	"gorm.io/gorm"
)

type BarbersRepository interface {
	AddBarber(req *models.Barber) error
	Update(id uint, barber models.Barber) error
	GetAll() ([]models.BarberResDTO, error)
	GetBarberByID(id uint) (*models.Barber, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type barbersRepository struct {
	db *gorm.DB
}

func NewBarbersRepository(db *gorm.DB) BarbersRepository {
	return &barbersRepository{db: db}
}

func (r *barbersRepository) AddBarber(req *models.Barber) error {
	if req == nil {
		return nil
	}
	return r.db.Create(req).Error
}

func (r *barbersRepository) Update(id uint, barber models.Barber) error {
	return r.db.Model(&models.Barber{}).Where("id = ?", id).Updates(barber).Error
}

func (r *barbersRepository) GetAll() ([]models.BarberResDTO, error) {
	var barbers []models.BarberResDTO

	if err := r.db.Model(&models.Barber{}).Find(&barbers).Error; err != nil {
		return nil, err
	}

	return barbers, nil
}

func (r *barbersRepository) GetBarberByID(id uint) (*models.Barber, error) {
	var barber models.Barber

	if err := r.db.First(&barber, id).Error; err != nil {
		return nil, err
	}

	return &barber, nil
}
func (r *barbersRepository) Delete(id uint) error {
	return r.db.Delete(&models.Barber{}, id).Error
}

func (r *barbersRepository) Exists(id uint) (bool, error){
	var count int64

	if err := r.db.Model(&models.Barber{}).Where("id = ?", id).Count(&count).Error; err != nil{
		return false, err
	}
	return count > 0, nil
}
