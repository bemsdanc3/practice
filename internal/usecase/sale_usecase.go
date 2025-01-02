package usecase

import (
	"errors"
	"tickets/internal/entities"
	"tickets/internal/repository"
	"time"
)

type SaleUsecase interface {
	ProcessSale(request entities.SaleRequest) error
}

type saleUsecase struct {
	SegmentRepo repository.SegmentRepository
}

func NewSaleUsecase(segmentRepo repository.SegmentRepository) SaleUsecase {
	return &saleUsecase{
		SegmentRepo: segmentRepo,
	}
}

func (u *saleUsecase) ProcessSale(request entities.SaleRequest) error {
	segments, err := u.SegmentRepo.FindSegmentsByTicketNumber(request.Passenger.TicketNumber)
	if err != nil {
		return err
	}
	if len(segments) > 0 {
		return errors.New("ticket already sold")
	}

	var newSegments []entities.Segment
	for i, route := range request.Routes {
		departDatetime := route.DepartDatetime.UTC()
		arriveDatetime := route.ArriveDatetime.UTC()

		_, offset := request.OperationTime.In(time.UTC).Zone()

		newSegment := entities.Segment{
			TicketNumber:        request.Passenger.TicketNumber,
			SegmentNumber:       i + 1, // Номер сегмента
			OperationType:       "sale",
			OperationTime:       request.OperationTime.UTC(),
			OperationTimeZone:   int16(offset / 3600), // Преобразуем смещение в часы
			OperationPlace:      request.OperationPlace,
			PassengerName:       request.Passenger.Name,
			PassengerSurname:    request.Passenger.Surname,
			PassengerPatronymic: request.Passenger.Patronymic,
			DocType:             request.Passenger.DocType,
			DocNumber:           request.Passenger.DocNumber,
			Birthdate:           time.Time{}, // Необходимо добавить дату рождения из запроса
			Gender:              request.Passenger.Gender,
			PassengerType:       request.Passenger.PassengerType,
			TicketType:          request.Passenger.TicketType,
			AirlineCode:         route.AirlineCode,
			FlightNum:           route.FlightNum,
			DepartPlace:         route.DepartPlace,
			DepartDatetime:      departDatetime,
			ArrivePlace:         route.ArrivePlace,
			ArriveDatetime:      arriveDatetime,
			PNRID:               route.PnrID,
			IsRefunded:          false, // Статус "не сдан"
		}
		newSegments = append(newSegments, newSegment)
	}

	if err := u.SegmentRepo.AddSaleSegment(newSegments); err != nil {
		return err
	}

	return nil
}
