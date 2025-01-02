package server

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"tickets/internal/entities"
	"tickets/internal/errs"
	"tickets/internal/usecase"
	validation "tickets/internal/validator"
)

type RefundHandler struct {
	RefundUsecase usecase.RefundUsecase
}

func NewRefundHandler(refundUsecase usecase.RefundUsecase) *RefundHandler {
	return &RefundHandler{
		RefundUsecase: refundUsecase,
	}
}

func (h *RefundHandler) HandleRefund(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Warn("Invalid request method", slog.String("method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 2<<10) // 2 KB
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Request body too large", slog.Any("error", err))
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	if err := validation.ValidateJSONWithSchemaFile(body, "internal/schemas/refund-schema.json"); err != nil {
		slog.Error("JSON validation failed", slog.Any("error", err))
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	var refundRequest entities.RefundRequest
	if err := json.Unmarshal(body, &refundRequest); err != nil {
		slog.Error("Invalid JSON format", slog.Any("error", err))
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	slog.Info("Refund request received",
		slog.String("ticketNumber", refundRequest.TicketNumber),
		slog.Time("operationTime", refundRequest.OperationTime),
		slog.String("operationPlace", refundRequest.OperationPlace),
	)

	err = h.RefundUsecase.RefundTicketWithContext(r.Context(), refundRequest.TicketNumber, refundRequest.OperationTime, refundRequest.OperationPlace)
	if err != nil {
		if errors.Is(err, errs.ErrTicketAlreadyRefunded) {
			slog.Warn("Ticket already refunded", slog.String("ticketNumber", refundRequest.TicketNumber))
			http.Error(w, "Ticket already refunded", http.StatusConflict)
			return
		}
		slog.Error("Error processing refund", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Refund processed successfully", slog.String("ticketNumber", refundRequest.TicketNumber))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Refund processed successfully"))
}
