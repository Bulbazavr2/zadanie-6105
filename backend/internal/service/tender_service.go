package service

import (
	"context"
	"tender_srevice/internal/domain"
	"tender_srevice/internal/repository"
	"fmt"
)

type TenderService struct {
	Repo *repository.PostgresRepository
}

func NewTenderService(repo *repository.PostgresRepository) *TenderService {
	return &TenderService{
		Repo: repo,
	}
}


type CreateTenderRequest struct {
	Name             string
	Description      string
	ServiceType      string
	Status           string
	OrganizationID   string
	CreatorUsername  string
}

func (s *TenderService) CreateTender(ctx context.Context, req CreateTenderRequest) (*domain.Tender, error) {
	newTender := &domain.Tender{
		Name:             req.Name,
		Description:      req.Description,
		ServiceType:      req.ServiceType,
		Status:           req.Status,
		OrganizationID:   req.OrganizationID,
		CreatorUsername:  req.CreatorUsername,
	}

	err := s.Repo.InsertTender(ctx, newTender)
	if err != nil {
		return nil, err
	}

	return newTender, nil
}

func (s *TenderService) GetTenders(ctx context.Context) ([]*domain.Tender, error) {
	return s.Repo.GetAllTenders(ctx)
}

func (s *TenderService) GetTendersByUsername(ctx context.Context, username string) ([]*domain.Tender, error) {
	return s.Repo.GetTendersByUsername(ctx, username)
}

func (s *TenderService) GetTenderStatus(ctx context.Context, tenderID, currentUsername string) (string, error) {
	tender, err := s.Repo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return "", err
	}

	if tender.Status == "PUBLISHED" {
		return tender.Status, nil
	}

	isResponsible, err := s.Repo.IsUserResponsibleForOrganization(ctx, currentUsername, tender.OrganizationID)
	if err != nil {
		fmt.Println("Ошибка при проверке ответственности:", err)
		return "", err
	}

	if !isResponsible {
		return "", fmt.Errorf("доступ запрещен: пользователь не является сотрудником организации")
	}

	return tender.Status, nil
}

type TenderUpdateRequest struct {
	Username 	*string `json:"username"`
	TenderID 	*string `json:"tenderId"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ServiceType *string `json:"serviceType"`
}

func (s *TenderService) UpdateTender(ctx context.Context, req TenderUpdateRequest) (*domain.Tender, error) {
	
	fmt.Println(req.Username)
	if req.TenderID == nil {
		return nil, fmt.Errorf("tender ID is required")
	}

	tender, err := s.Repo.GetTenderByID(ctx, *req.TenderID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Creator Username:", tender.CreatorUsername)
	
	if req.Username == nil {
		return nil, fmt.Errorf("username is required service")
	}
	fmt.Println("Request Username:", *req.Username)

	if tender.CreatorUsername != *req.Username {
		return nil, fmt.Errorf("unauthorized")
	}

	if req.Name != nil {
		nameValue := *req.Name
		fmt.Println("Name value:", nameValue)
		tender.Name = nameValue
	}

	if req.Description != nil {
		tender.Description = *req.Description
	}
	if req.ServiceType != nil {
		tender.ServiceType = *req.ServiceType
	}

	err = s.Repo.UpdateTender(ctx, tender)
	if err != nil {
		return nil, err
	}

	return tender, nil
}

func (s *TenderService) UpdateTenderStatus(ctx context.Context, tenderID, newStatus, currentUsername string) (*domain.Tender, error) {
	
	
	tender, err := s.Repo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// if tender.CreatorUsername != currentUsername {
	// 	return nil, fmt.Errorf("unauthorized")
	// }

	isResponsible, err := s.Repo.IsUserResponsibleForOrganization(ctx, currentUsername, tender.OrganizationID)
    if err != nil {
		fmt.Println("Errordasd:", err)
        return nil, err
    }

	if !isResponsible {
        return nil, fmt.Errorf("unauthorized")
    }

	tender.Status = newStatus  // Добавьте эту строку

	err = s.Repo.UpdateTenderStatus(ctx, tender)
	if err != nil {
		return nil, err
	}

	return tender, nil
}

func (s *TenderService) RollbackTender(ctx context.Context, tenderID string, version int, username string) (*domain.Tender, error) {
	tender, err := s.Repo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	isResponsible, err := s.Repo.IsUserResponsibleForOrganization(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	if !isResponsible {
		return nil, fmt.Errorf("unauthorized")
	}

	updatedTender, err := s.Repo.RollbackTender(ctx, tenderID, version)
	if err != nil {
		return nil, err
	}

	return updatedTender, nil
}

