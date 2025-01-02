package app

import (
	"log/slog"
	"net/http"
	"os"
	"tickets/internal/delivery/server"
	"tickets/internal/repository"
	"tickets/internal/usecase"
	"tickets/pkg/db"
)

const (
	ServerPort        = ":8080"
	ProcessSaleURL    = "/process/sale/"
	ProcessRefundURL  = "/process/refund/"
	LogInitialization = "Initializing database..."
	LogServerStart    = "Starting server..."
	LogServerError    = "Error starting server"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info(LogInitialization)
	db.InitDB()

	// определение репозиториев
	segmentRepository := repository.NewSegmentRepository(db.GetDB())
	refundRepository := repository.NewRefundRepository(db.GetDB())

	// определение юзкейсов
	saleUsecase := usecase.NewSaleUsecase(segmentRepository)
	refundUsecase := usecase.NewRefundUsecase(refundRepository)

	// определение обработчиков
	segmentHandler := server.NewSaleHandler(saleUsecase)
	refundHandler := server.NewRefundHandler(refundUsecase)

	// урлы
	http.HandleFunc(ProcessSaleURL, segmentHandler.HandleSale)
	http.HandleFunc(ProcessRefundURL, refundHandler.HandleRefund)

	slog.Info(LogServerStart, slog.String("port", ServerPort))
	err := http.ListenAndServe(ServerPort, nil)
	if err != nil {
		slog.Error(LogServerError, slog.Any("error", err))
	}
}
