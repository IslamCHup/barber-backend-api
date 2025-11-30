package service

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
)

type ClientService interface {
	AddClient(req *models.ClientCreateReqDTO) error
	GetClientByID(id uint) (*models.Client, error)
	GetAllClients() ([]models.ClientRespDTO, error)
	Update(id uint, client models.ClientUpdateReqDTO) error
	Delete(id uint) error
}

type clientService struct {
	service repository.ClientsRepository
}

func NewClientsService(service repository.ClientsRepository) ClientService {
	return &clientService{service: service}
}

func (s *clientService) AddClient(req *models.ClientCreateReqDTO) error {
	res := models.Client{
		FullName: req.FullName,
	}
	if err := s.service.AddClient(&res); err != nil {
		return err
	}
	return nil
}

func (s *clientService) GetClientByID(id uint) (*models.Client, error) {
	return s.service.GetClientByID(id)
}

func (s *clientService) GetAllClients() ([]models.ClientRespDTO, error) {
	return s.service.GetAllClients()
}

func (s *clientService) Update(id uint, client models.ClientUpdateReqDTO) error {
	if err := s.service.Update(id, client); err != nil {
		return err
	}
	return nil
}

func (s *clientService) Delete(id uint) error {
	client, err := s.service.GetClientByID(id)
	if err != nil {
		return err
	}
	return s.service.Delete(client.ID)
}
