package usecase

import (
	"context"
	"tickets/internal/errs"
	"tickets/internal/repository"
	"time"
)

type RefundUsecase interface {
	RefundTicketWithContext(ctx context.Context, ticketNumber string, operationTime time.Time, operationPlace string) error
}

type refundUsecase struct {
	RefundRepo repository.RefundRepository
}

func NewRefundUsecase(refundRepo repository.RefundRepository) RefundUsecase {
	return &refundUsecase{
		RefundRepo: refundRepo,
	}
}

func (u *refundUsecase) RefundTicketWithContext(ctx context.Context, ticketNumber string, operationTime time.Time, operationPlace string) error {
	// проверяем статус билета перед выполнением операции
	status, err := u.RefundRepo.GetTicketStatus(ctx, ticketNumber)
	if err != nil {
		return errs.ErrTicketNotFound
	}

	switch status {
	case "refunded":
		return errs.ErrTicketAlreadyRefunded
	case "sold":
		return errs.ErrTicketAlreadySold
	}

	return u.RefundRepo.MarkSegmentsAsRefunded(ticketNumber, operationTime, operationPlace)
}
