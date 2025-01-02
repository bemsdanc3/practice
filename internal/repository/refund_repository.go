package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tickets/internal/entities"
	"tickets/internal/errs"
	"time"
)

type RefundRepository interface {
	MarkSegmentsAsRefunded(ticketNumber string, operationTime time.Time, operationPlace string) error
	GetTicketStatus(ctx context.Context, ticketNumber string) (string, error)
}

type refundRepository struct {
	DB *gorm.DB
}

func NewRefundRepository(db *gorm.DB) RefundRepository {
	return &refundRepository{
		DB: db,
	}
}

func (r *refundRepository) MarkSegmentsAsRefunded(ticketNumber string, operationTime time.Time, operationPlace string) error {
	tx := r.DB.Begin()

	var segments []entities.Segment
	err := r.DB.Where("ticket_number = ?", ticketNumber).Find(&segments).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(segments) == 0 {
		tx.Rollback()
		return errors.New("ticket not found")
	}

	// проверка сдан ли уже билет
	for _, segment := range segments {
		if segment.OperationType == "refund" {
			tx.Rollback()
			return errors.New("ticket already refunded")
		}
	}

	for _, segment := range segments {
		segment.IsRefunded = true
		segment.OperationType = "refund"
		segment.OperationTime = operationTime
		segment.OperationPlace = operationPlace

		if err := r.DB.Save(&segment).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *refundRepository) GetTicketStatus(ctx context.Context, ticketNumber string) (string, error) {
	var segment entities.Segment
	err := r.DB.WithContext(ctx).Where("ticket_number = ?", ticketNumber).First(&segment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errs.ErrTicketNotFound
	}
	return segment.OperationType, nil
}
