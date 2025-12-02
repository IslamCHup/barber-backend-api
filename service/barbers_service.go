package service

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
	"errors"
	"log/slog"
)

type BarberService interface {
	AddBarber(req *models.BarbersCreateReqDTO) error
	Update(id uint, barberInp models.Barber) (*models.Barber, error)
	GetAll() ([]models.BarberResDTO, error)
	GetBarberByID(id uint) (*models.Barber, error)
	Delete(id uint) error
}

type barberService struct {
	logger  *slog.Logger
	service repository.BarbersRepository
}

func NewBarbersService(logger *slog.Logger, service repository.BarbersRepository) BarberService {
	return &barberService{logger: logger, service: service}
}

func (s *barberService) AddBarber(req *models.BarbersCreateReqDTO) error {
	s.logger.Info("запуск функции добавления записи парикмахер",
		"req", req.FullName,
	)

	if req.FullName == "" {
		s.logger.Error("передан пустой запрос",
			"op", "service.barber.AddBarber",
			"error", "empty req")
		return errors.New("empty req")
	}
	barber := models.Barber{
		FullName: req.FullName,
	}

	s.logger.Info("присвоено значение из DTO в доменную модель",
		"доменная модель", barber.FullName,
	)

	if err := s.service.AddBarber(&barber); err != nil {
		s.logger.Error("ошибка добавления парикмахера",
			"op", "service.barber.AddBarber",
			"error", err,
		)
		return err
	}
	s.logger.Info("добавление записи прошло успешно",
		"op", "service.barber.AddBarber",
	)
	return nil
}

func (s *barberService) Update(id uint, barberInp models.Barber) (*models.Barber, error) {
	s.logger.Info("запуск функции обновления парикмахера", "op", "service.barber.Update", "id", id)

	isExist, err := s.service.Exists(id)
	if err != nil {
		s.logger.Error("ошибка проверки существования записи",
			"op", "service.barber.Update",
			"id", id,
			"error", err,
		)
		return nil, err
	}
	if !isExist {
		s.logger.Error("записи с таким ID не найдено",
			"op", "service.barber.Update",
			"id", id,
		)
		return nil, errors.New("записи с таким ID не найдено")
	}

	if err := s.service.Update(id, barberInp); err != nil {
		s.logger.Error("ошибка обновления записи парикмахера",
			"op", "service.barber.Update",
			"id", id,
			"error", err,
		)
		return nil, err
	}
	s.logger.Info("обновление прошло успешно",
		"op", "service.barber.Update",
		"id", id,
	)

	barber, err := s.service.GetBarberByID(id)
	if err != nil {
		s.logger.Error("ошибка получения парикмахера после обновления",
			"op", "service.barber.Update",
			"id", id,
			"error", err,
		)
		return nil, errors.New("id не найден")
	}

	s.logger.Info("функция обновления завершена успешно",
		"op", "service.barber.Update",
		"id", barber.ID,
		"full_name", barber.FullName,
	)

	return barber, nil
}

func (s *barberService) GetAll() ([]models.BarberResDTO, error) {
	s.logger.Info("запуск функции получения всех парикмахеров", "op", "service.barber.GetAll")
	res, err := s.service.GetAll()
	if err != nil {
		s.logger.Error("ошибка получения списка парикмахеров",
			"op", "service.barber.GetAll",
			"error", err,
		)
		return nil, err
	}
	s.logger.Info("список парикмахеров получен успешно",
		"op", "service.barber.GetAll",
		"count", len(res),
	)
	return res, nil
}

func (s *barberService) GetBarberByID(id uint) (*models.Barber, error) {
	s.logger.Info("запуск функции получения парикмахера по ID",
		"op", "service.barber.GetBarberByID",
		"id", id,
	)
	barber, err := s.service.GetBarberByID(id)
	if err != nil {
		s.logger.Error("ошибка запроса парикмахера",
			"op", "service.barber.GetBarberByID",
			"id", id,
			"error", err,
		)
		return nil, err
	}
	if barber == nil {
		s.logger.Error("запись не найдена",
			"op", "service.barber.GetBarberByID",
			"id", id,
		)
		return nil, errors.New("запись не найдена")
	}
	s.logger.Info("парикмахер по ID найден",
		"op", "service.barber.GetBarberByID",
		"id", barber.ID,
		"full_name", barber.FullName,
	)

	return barber, nil
}

func (s *barberService) Delete(id uint) error {
	s.logger.Info("запуск функции удаления парикмахера",
		"op", "service.barber.Delete",
		"id", id,
	)
	barber, err := s.service.GetBarberByID(id)
	if err != nil {
		s.logger.Error("ошибка получения парикмахера перед удалением",
			"op", "service.barber.Delete",
			"id", id,
			"error", err,
		)
		return err
	}
	if barber == nil {
		s.logger.Error("запись не найдена для удаления",
			"op", "service.barber.Delete",
			"id", id,
		)
		return errors.New("запись не найдена")
	}

	if err := s.service.Delete(barber.ID); err != nil {
		s.logger.Error("ошибка удаления парикмахера",
			"op", "service.barber.Delete",
			"id", barber.ID,
			"error", err,
		)
		return err
	}

	s.logger.Info("удаление прошло успешно",
		"op", "service.barber.Delete",
		"id", barber.ID,
		"full_name", barber.FullName,
	)

	return nil
}
