package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tender_srevice/internal/service"

	"github.com/gorilla/mux"
)

type BidHandler struct {
	service *service.BidService
}

func NewBidHandler(service *service.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var req service.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bid, err := h.service.CreateBid(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bid)
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	var req struct {
		Username   string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetBidStatus(r.Context(), bidID, req.Username)
	if err != nil {
		http.Error(w, "Failed to get bid status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *BidHandler) GetMyBids(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Здесь нужно получить ID пользователя по его username
	userID, err := h.service.Repo.GetUserIDByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	bids, err := h.service.GetBidsByAuthorID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bids)
}

func (h *BidHandler) GetBidsByTenderID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	bids, err := h.service.GetBidsByTenderID(r.Context(), tenderID, req.Username)
	if err != nil {
		if err.Error() == "user is not authorized to view bids for this tender" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			log.Printf("Error getting bids: %v", err) // Добавьте эту строку для логирования
			http.Error(w, fmt.Sprintf("Failed to get bids: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bids)
}

func (h *BidHandler) GetBidStatusByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetBidStatus(r.Context(), bidID, req.Username)
	if err != nil {
		if err.Error() == "user is not authorized to view this bid status" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get bid status: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *BidHandler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	var req struct {
		Username string `json:"username"`
		NewStatus string `json:"newStatus"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.NewStatus == "" {
		http.Error(w, "Требуется указать имя пользователя и новый статус", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateBidStatus(r.Context(), bidID, req.Username, req.NewStatus)
	if err != nil {
		if err != nil {
			if err.Error() == "user is not authorized to view this bid status" {
				http.Error(w, err.Error(), http.StatusForbidden)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get bid status: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Статус заявки успешно обновлен"})
}

// EditBidRequest представляет структуру для запроса на редактирование заявки
type EditBidRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var req EditBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.EditBid(r.Context(), bidID, username, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid updated successfully"})
}