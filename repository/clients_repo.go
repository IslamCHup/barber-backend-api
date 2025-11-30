package repository

import (
	"barber-backend-api/internal/models"

	"gorm.io/gorm"
)

type ClientsRepository interface {
	AddClient(req *models.Client) error
	GetClientByID(id uint) (*models.Client, error)
	GetAllClients() ([]models.ClientRespDTO, error)
	Update(id uint, client models.ClientUpdateReqDTO) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type clientsRepository struct {
	db *gorm.DB
}

func NewClientsRepository(db *gorm.DB) ClientsRepository {
	return &clientsRepository{db: db}
}

func (r *clientsRepository) AddClient(req *models.Client) error {
	if req == nil {
		return nil
	}
	return r.db.Create(req).Error
}

func (r *clientsRepository) GetClientByID(id uint) (*models.Client, error) {
	var client models.Client

	if err := r.db.First(&client, id).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *clientsRepository) GetAllClients() ([]models.ClientRespDTO, error) {
	var clients []models.ClientRespDTO

	if err := r.db.Model(&models.Client{}).Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientsRepository) Update(id uint, client models.ClientUpdateReqDTO) error {
	return r.db.Model(&models.Client{}).Where("id = ?", id).Updates(client).Error
}

func (r *clientsRepository) Delete(id uint) error {
	return r.db.Delete(&models.Client{}, id).Error
}

func (r *clientsRepository) Exists(id uint) (bool, error){
	var count int64

	if err := r.db.Model(&models.Client{}).Where("id = ?", id).Count(&count).Error; err != nil{
		return false, err
	}
	return count > 0, nil
}