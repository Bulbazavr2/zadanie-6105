package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"log"

	"github.com/gorilla/mux"

	"tender_srevice/internal/domain"
	"tender_srevice/internal/service"
)



type TenderHandler struct {
	service *service.TenderService
}


func NewTenderHandler(service *service.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	// var req struct {
	// 	Name           string `json:"name"`
	// 	Description    string `json:"description"`
	// 	ServiceType    string `json:"serviceType"`
	// 	Status         string `json:"status"`
	// 	OrganizationID string  `json:"organizationId"`
	// 	CreatorUsername string `json:"creatorUsername"`
	// }

	var req domain.Tender

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tender, err := h.service.CreateTender(r.Context(), service.CreateTenderRequest{
		Name:             req.Name,
		Description:      req.Description,
		ServiceType:      req.ServiceType,
		Status:           req.Status,
		OrganizationID:   req.OrganizationID,
		CreatorUsername:  req.CreatorUsername,
	})

	if err != nil {
		http.Error(w, "Failed to create tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {

	
	tenders, err := h.service.GetTenders(r.Context())
	if err != nil {
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func (h *TenderHandler) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required handler", http.StatusBadRequest)
		return
	}

	tenders, err := h.service.GetTendersByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]
	username := r.URL.Query().Get("username")

	if tenderID == "" {
		http.Error(w, "Tender ID is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetTenderStatus(r.Context(), tenderID, username)
	if err != nil {
		if err.Error() == "tender not found" {
			http.Error(w, "Tender not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get tender status", http.StatusInternalServerError)
		}
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *TenderHandler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	var req struct {
		Status   string `json:"status"`
		Username string `json:"username"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка при декодировании запроса: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Status != domain.TenderStatusCreated && 
	   req.Status != domain.TenderStatusPublished && 
	   req.Status != domain.TenderStatusClosed {
		http.Error(w, "Недопустимый статус тендера", http.StatusBadRequest)
		return
	}

	fmt.Println(tenderID, req.Status, req.Username)
	updatedTender, err := h.service.UpdateTenderStatus(r.Context(), tenderID, req.Status, req.Username)
	if err != nil {
		if err.Error() == "tender not found" {
			http.Error(w, "Тендер не найден", http.StatusNotFound)
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Нет прав для изменения статуса тендера", http.StatusForbidden)
		} else {
			http.Error(w, "Ошибка при обновлении статуса", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTender)
}



type TenderUpdateRequest struct {
	Username    string  `json:"username"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ServiceType *string `json:"serviceType"`
}

func (h *TenderHandler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	var req TenderUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	

	updatedTender, err := h.service.UpdateTender(r.Context(), service.TenderUpdateRequest{
		Username: &req.Username,
		TenderID: &tenderID,
		Name: req.Name,
		Description: req.Description,
		ServiceType: req.ServiceType,
	})
	if err != nil {
		log.Printf("Ошибка при обновлении тендера: %v", err)
		if err.Error() == "tender not found" {
			http.Error(w, "Тендер не найден", http.StatusNotFound)
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Нет прав для изменения тендера", http.StatusForbidden)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при обновлении тендера: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTender)
}

func (h *TenderHandler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]
	versionStr := vars["version"]

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if tenderID == "" || versionStr == "" || req.Username == "" {
		http.Error(w, "Необходимы ID тендера, версия и имя пользователя", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Неверный формат версии", http.StatusBadRequest)
		return
	}

	updatedTender, err := h.service.RollbackTender(r.Context(), tenderID, version, req.Username)
	if err != nil {
		switch err.Error() {
		case "tender not found":
			http.Error(w, "Тендер не найден", http.StatusNotFound)
		case "unauthorized":
			http.Error(w, "Нет прав для отката тендера", http.StatusForbidden)
		case "version not found":
			http.Error(w, "Указанная версия не найдена", http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Ошибка при откате тендера: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTender)
}


