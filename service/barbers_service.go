package service

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
	"errors"
)

type BarberService interface {
	AddBarber(req *models.BarbersCreateReqDTO) error
	Update(id uint, barberInp models.Barber) (*models.Barber, error)
	GetAll() ([]models.BarberResDTO, error)
	GetBarberByID(id uint) (*models.Barber, error)
	Delete(id uint) error
}

type barberService struct {
	service repository.BarbersRepository
}

func NewBarbersService(service repository.BarbersRepository) BarberService {
	return &barberService{service: service}
}

func (s *barberService) AddBarber(req *models.BarbersCreateReqDTO) error {
	
	barber := models.Barber{
		FullName:  req.FullName,
	}
	if err := s.service.AddBarber(&barber); err != nil {
		return err
	}

	return nil
}

func (s *barberService) Update(id uint, barberInp models.Barber) (*models.Barber, error) {
	isExist, err := s.service.Exists(id)
	if err != nil || !isExist {
		return nil, errors.New("записи с таким ID не найдено")
	}

	if err := s.service.Update(id, barberInp); err != nil {
		return nil, err
	}

	barber, err := s.service.GetBarberByID(id)
	if err != nil {
		return nil, errors.New("id не найден")
	}

	return barber, err
}

func (s *barberService) GetAll() ([]models.BarberResDTO, error) {
	return s.service.GetAll()
}

func (s *barberService) GetBarberByID(id uint) (*models.Barber, error) {
	barber, err :=  s.service.GetBarberByID(id)
	if err != nil{
		return nil, err
	}
	if barber == nil{
		return nil, errors.New("запись не найдена")
	}
	return barber, nil
}

func (s *barberService) Delete(id uint) error {
	barber, err := s.service.GetBarberByID(id)
	if err != nil {
		return err
	}
	return s.service.Delete(barber.ID)
}
