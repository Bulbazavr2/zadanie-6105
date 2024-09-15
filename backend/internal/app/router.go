package app

import (
	"net/http"
	"tender_srevice/internal/config"
	"tender_srevice/internal/handler"
	"tender_srevice/internal/service"
	"tender_srevice/internal/repository"

	"github.com/gorilla/mux"
)

func SetupRouter(cfg *config.Config, repo *repository.PostgresRepository) *mux.Router {
	router := mux.NewRouter()

	tenderService := service.NewTenderService(repo)
	tenderHandler := handler.NewTenderHandler(tenderService)

	router.HandleFunc("/api/ping", handler.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/new", tenderHandler.CreateTender).Methods(http.MethodPost)
	router.HandleFunc("/api/tenders", tenderHandler.GetTenders).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/my", tenderHandler.GetMyTenders).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/{tenderId}/status", tenderHandler.GetTenderStatus).Methods(http.MethodGet)
	router.HandleFunc("/api/tenders/{tenderId}/status", tenderHandler.UpdateTenderStatus).Methods(http.MethodPut)
	router.HandleFunc("/api/tenders/{tenderId}/edit", tenderHandler.UpdateTender).Methods(http.MethodPatch)
	router.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender).Methods(http.MethodPut)

	bidService := service.NewBidService(repo)
	bidHandler := handler.NewBidHandler(bidService)

	router.HandleFunc("/api/bids/new", bidHandler.CreateBid).Methods(http.MethodPost)
	router.HandleFunc("/api/bids/my", bidHandler.GetMyBids).Methods(http.MethodGet)
	router.HandleFunc("/api/bids/{tenderId}/list", bidHandler.GetBidsByTenderID).Methods(http.MethodGet)
	router.HandleFunc("/api/bids/{bidId}/status", bidHandler.GetBidStatusByID).Methods(http.MethodGet)
	router.HandleFunc("/api/bids/{bidId}/status", bidHandler.UpdateBidStatus).Methods(http.MethodPut)
	router.HandleFunc("/api/bids/{bidId}/edit", bidHandler.EditBid).Methods(http.MethodPatch)


	router.HandleFunc("/api/tenders/{tenderId}/bids", bidHandler.GetBidsByTenderID).Methods(http.MethodGet)
	

	return router
}