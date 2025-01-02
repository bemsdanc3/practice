package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"tickets/internal/entities"
	"tickets/internal/usecase"
	validation "tickets/internal/validator"
)

type SaleHandler struct {
	SaleUsecase usecase.SaleUsecase
}

func NewSaleHandler(saleUsecase usecase.SaleUsecase) *SaleHandler {
	return &SaleHandler{
		SaleUsecase: saleUsecase,
	}
}

func (h *SaleHandler) HandleSale(w http.ResponseWriter, r *http.Request) {
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

	if err := validation.ValidateJSONWithSchemaFile(body, "internal/schemas/sale-schema.json"); err != nil {
		slog.Error("JSON validation failed", slog.Any("error", err))
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	var saleRequest entities.SaleRequest
	if err := json.Unmarshal(body, &saleRequest); err != nil {
		slog.Error("Invalid JSON format", slog.Any("error", err))
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	slog.Info("Sale request received", slog.String("ticketNumber", saleRequest.Passenger.TicketNumber))

	if err := h.SaleUsecase.ProcessSale(saleRequest); err != nil {
		if err.Error() == "ticket already sold" {
			slog.Warn("Ticket already sold", slog.String("ticketNumber", saleRequest.Passenger.TicketNumber))
			http.Error(w, "Ticket already sold", http.StatusConflict)
			return
		}
		slog.Error("Error processing sale", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Sale processed successfully", slog.String("ticketNumber", saleRequest.Passenger.TicketNumber))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sale processed successfully"))
}
