package repository

import (
	"barber-backend-api/internal/models"
	"errors"
	"log/slog"

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
	logger *slog.Logger
	db     *gorm.DB
}

func NewBarbersRepository(logger *slog.Logger, db *gorm.DB) BarbersRepository {
	return &barbersRepository{logger: logger, db: db}
}

func (r *barbersRepository) AddBarber(req *models.Barber) error {
	r.logger.Debug("начало добавления записи парикмахера",
		"op", "repo.barber.create")

	if req == nil {
		err := errors.New("передан пустой запрос")
		r.logger.Error("ошибка добавления парикмахера",
			"op", "repo.barber.create",
			"error", err,
		)
		return err
	}

	if err := r.db.Create(req).Error; err != nil {
		r.logger.Error("не удалось добавить парикмахера",
			"op", "repo.barber.create",
			"error", err,
			"name", req.FullName,
		)
		return err
	}

	r.logger.Debug("парикмахер успешно добавлен",
		"op", "repo.barber.create",
		"id", req.ID,
		"name", req.FullName,
	)
	return nil
}

func (r *barbersRepository) Update(id uint, barber models.Barber) error {
	r.logger.Debug("старт обновления записи парикмахера",
		"op", "repo.barber.update",
		"id", id,
	)

	result := r.db.Model(&models.Barber{}).Where("id = ?", id).Updates(barber)
	if result.Error != nil {
		r.logger.Error("не удалось обновить запись парикмахера",
			"op", "repo.barber.update",
			"id", id,
			"error", result.Error,
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("запись для обновления не найдена",
			"op", "repo.barber.update",
			"id", id,
		)
		return gorm.ErrRecordNotFound
	}

	r.logger.Debug("обновление записи прошло успешно",
		"op", "repo.barber.update",
		"id", id,
		"rows_affected", result.RowsAffected,
	)
	return nil
}

func (r *barbersRepository) GetAll() ([]models.BarberResDTO, error) {
	r.logger.Debug("запуск функции для получения всех парикмахеров из БД",
		"op", "repo.barber.get_all",
	)

	var barbers []models.BarberResDTO
	result := r.db.Model(&models.Barber{}).Find(&barbers)
	if result.Error != nil {
		r.logger.Error("не удалось найти записи парикмахеров",
			"op", "repo.barber.get_all",
			"error", result.Error,
		)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("таблица barber пустая",
			"op", "repo.barber.get_all",
		)
		return []models.BarberResDTO{}, nil
	}

	r.logger.Debug("получение всех записей barber прошло успешно",
		"op", "repo.barber.get_all",
		"rows", len(barbers),
	)
	return barbers, nil
}

func (r *barbersRepository) GetBarberByID(id uint) (*models.Barber, error) {
	r.logger.Debug("получение barber по ID",
		"op", "repo.barber.get_by_id",
		"id", id,
	)

	var barber models.Barber
	result := r.db.First(&barber, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Warn("barber не найден",
				"op", "repo.barber.get_by_id",
				"id", id,
			)
			return nil, result.Error
		}
		r.logger.Error("ошибка при получении barber по ID",
			"op", "repo.barber.get_by_id",
			"id", id,
			"error", result.Error,
		)
		return nil, result.Error
	}

	r.logger.Debug("получение barber по ID прошло успешно",
		"op", "repo.barber.get_by_id",
		"id", barber.ID,
		"barber", barber.FullName,
	)
	return &barber, nil
}

func (r *barbersRepository) Delete(id uint) error {
	r.logger.Debug("начало удаления записи по ID",
		"op", "repo.barber.delete",
		"id", id,
	)

	result := r.db.Delete(&models.Barber{}, id)
	if result.Error != nil {
		r.logger.Error("не удалось удалить запись по ID",
			"op", "repo.barber.delete",
			"id", id,
			"error", result.Error,
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("нет записей для удаления",
			"op", "repo.barber.delete",
			"id", id,
		)
		return gorm.ErrRecordNotFound
	}

	r.logger.Debug("запись успешно удалена",
		"op", "repo.barber.delete",
		"id", id,
		"rows_affected", result.RowsAffected,
	)
	return nil
}

func (r *barbersRepository) Exists(id uint) (bool, error) {
	r.logger.Debug("проверка на существование записи по ID запущена",
		"op", "repo.barber.exists",
		"id", id,
	)

	var count int64
	result := r.db.Model(&models.Barber{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		r.logger.Error("не удалось проверить запись на существование",
			"op", "repo.barber.exists",
			"id", id,
			"error", result.Error,
		)
		return false, result.Error
	}

	isExist := count > 0
	r.logger.Debug("проверка на существование записи по ID прошла успешно",
		"op", "repo.barber.exists",
		"id", id,
		"exists", isExist,
	)
	return isExist, nil
}
