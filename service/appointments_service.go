package service

import (
	"barber-backend-api/internal/models"
	"barber-backend-api/repository"
	"errors"
	"time"
)

type AppointmentsService interface {
	GetAllAppointments() ([]models.Appointments, error)
	CreateAppointment(req *models.AppointmentsCreateDTO) error
	Update(barberID uint, req models.AppointmentsUpdateReqDTO) error
	Delete(id uint) error
	GetAllAppointmentsByBarberID(id uint) ([]models.Appointments, error)
	GetByID(id uint) (*models.Appointments, error)
}

type appointmentsService struct {
	service repository.AppointmentsRepository
	barber  repository.BarbersRepository
}

func NewAppointmentsService(service repository.AppointmentsRepository, barber repository.BarbersRepository) AppointmentsService {
	return &appointmentsService{service: service, barber: barber}
}

func (s *appointmentsService) GetAllAppointments() ([]models.Appointments, error) {
	return s.service.GetAllAppointments()
}

func (s *appointmentsService) CreateAppointment(req *models.AppointmentsCreateDTO) error {
	listAppmts, err := s.service.GetAllAppointmentsByBarberID(req.BarberID)
	if err != nil {
		return err
	}
	date := req.Time
	t, err := time.Parse("2006-01-02 15", date)
	if err != nil {
		return errors.New("неправильный формат времени, нужен YYYY-MM-DD HH")
	}

	if t.Minute() != 0 || t.Second() != 0 {
		t = t.Add(time.Hour).Truncate(time.Hour)
	}

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return errors.New("нельзя записаться на выходной день")
	}

	if t.Before(time.Now()) {
		return errors.New("запись просрочена")
	}

	if t.Hour() < 9 || t.Hour() >= 17 {
		return errors.New("парикмахер работает с 9 до 17 часов")
	}

	for _, v := range listAppmts {
		vt, err := time.Parse(time.DateTime, v.Time)
		if err != nil {
			continue
		}
		existingEnd := vt.Add(time.Hour)
		newEnd := t.Add(time.Hour)
		if vt.Before(newEnd) && existingEnd.After(t) {
			return errors.New("это время занято другими")
		}
	}
	req.Time = t.Format(time.DateTime)

	requestInput := models.Appointments{
		BarberID: req.BarberID,
		ClientID: req.ClientID,
		Time:     req.Time,
	}

	if err := s.service.CreateAppointment(&requestInput); err != nil {
		return err
	}
	return nil
}

func (s *appointmentsService) Update(id uint, req models.AppointmentsUpdateReqDTO) error {
	lastAppointments, err := s.service.GetLastAppointments(req.BarberID)
	if err != nil {
		return err
	}

	if *lastAppointments.Rating != 0{
		return errors.New("вы уже ставили оценку")
	}

	startVisit, _ := time.Parse("2006-01-02 15", lastAppointments.Time)
	endVisit := startVisit.Add(time.Hour)

	if endVisit.After(time.Now()) {
		return errors.New("оценку можно ставить после услуги")
	}

	if err := s.service.Update(lastAppointments.ID, req); err != nil {
		return err
	}

	if err := s.ratingBarbers(req.BarberID); err != nil {
		return err
	}

	return nil
}

func (s *appointmentsService) Delete(id uint) error {
	appointment, err := s.service.GetByID(id)
	if err != nil {
		return err
	}

	startVisit, _ := time.Parse("2006-01-02 15", appointment.Time)
	endVisit := startVisit.Add(time.Hour)

	if endVisit.After(time.Now()) {
		return errors.New("невозможно удалить запись после предоставления услуги ")
	}

	if err := s.service.Delete(appointment.ID); err != nil {
		return err
	}
	return nil
}

func (s *appointmentsService) GetAllAppointmentsByBarberID(id uint) ([]models.Appointments, error) {
	return s.service.GetAllAppointmentsByBarberID(id)
}

func (s *appointmentsService) GetByID(id uint) (*models.Appointments, error) {
	return s.service.GetByID(id)
}

func (s *appointmentsService) ratingBarbers(barberID uint) error {
	appointments, err := s.service.GetAllAppointmentsByBarberID(barberID)
	if err != nil {
		return err
	}

	sum := 0
	count := 0

	for _, v := range appointments {
		if v.Rating != nil {
			sum += *v.Rating
			count++
		}
	}

	avg_rating := 0.0
	if count == 0 {
		avg_rating = 0
	} else {
		avg_rating = float64(sum) / float64(count)
	}
	ratingInput := models.Barber{
		AvgRating: avg_rating,
	}

	s.barber.Update(barberID, ratingInput)
	return nil
}
