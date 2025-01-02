package repository

import (
	"gorm.io/gorm"
	"tickets/internal/entities"
)

type SegmentRepository interface {
	AddSaleSegment(segments []entities.Segment) error
	FindSegmentsByTicketNumber(ticketNumber string) ([]entities.Segment, error)
}

type segmentRepository struct {
	DB *gorm.DB
}

func NewSegmentRepository(db *gorm.DB) SegmentRepository {
	return &segmentRepository{
		DB: db,
	}
}

func (r *segmentRepository) AddSaleSegment(segments []entities.Segment) error {
	tx := r.DB.Begin()

	for _, segment := range segments {
		if err := tx.Create(&segment).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *segmentRepository) FindSegmentsByTicketNumber(ticketNumber string) ([]entities.Segment, error) {
	var segments []entities.Segment
	err := r.DB.Where("ticket_number = ? AND is_refunded = false", ticketNumber).Find(&segments).Error
	if err != nil {
		return nil, err
	}
	return segments, nil
}
