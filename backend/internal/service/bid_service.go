package service

import (
	"context"
	"fmt"
	"tender_srevice/internal/domain"
	"tender_srevice/internal/repository"
	
)

type BidService struct {
	Repo *repository.PostgresRepository
}

func NewBidService(repo *repository.PostgresRepository) *BidService {
	return &BidService{Repo: repo}
}

type CreateBidRequest struct {
	Name        string
	Description string
	TenderID    string
	AuthorType  string
	AuthorID    string
}

func (s *BidService) CreateBid(ctx context.Context, req CreateBidRequest) (*domain.Bid, error) {
	newBid := &domain.Bid{
		Name:        req.Name,
		Description: req.Description,
		Status:      domain.BidStatusPending,
		TenderID:    req.TenderID,
		AuthorType:  req.AuthorType,
		AuthorID:    req.AuthorID,
	}

	err := s.Repo.InsertBid(ctx, newBid)
	if err != nil {
		return nil, err
	}

	return newBid, nil
}

func (s *BidService) GetBidsByTenderID(ctx context.Context, tenderID, username string) ([]*domain.Bid, error) {
	// Проверяем, принадлежит ли пользователь к организации тендера
	
	isAllowed, err := s.Repo.IsUserInTenderOrganization(ctx, username, tenderID)
	if err != nil {
		fmt.Println("Error checking user permissions:", err)
		return nil, fmt.Errorf("failed to check user permissions: %w", err)

	}
	if !isAllowed {
		return nil, fmt.Errorf("user is not authorized to view bids for this tender")
	}

	return s.Repo.GetBidsByTenderID(ctx, tenderID)
}

func (s *BidService) UpdateBidStatus(ctx context.Context, bidID, username, newStatus string) error {
	// Проверяем, существует ли заявка
	bid, err := s.Repo.GetBidByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("failed to get bid: %w", err)
	}

	// Проверяем, принадлежит ли пользователь к организации тендера
	isAllowed, err := s.Repo.IsUserInTenderOrganization(ctx, username, bid.TenderID)
	if err != nil {
		fmt.Println("Error checking user permissions:", err)
		return fmt.Errorf("failed to check user permissions: %w", err)

	}
	if !isAllowed {
		return fmt.Errorf("user is not authorized to view bids for this tender")
	}

	// Обновляем статус заявки
	err = s.Repo.UpdateBidStatus(ctx, bidID, newStatus)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса заявки: %w", err)
	}

	return nil
}

func (s *BidService) GetBidsByAuthorID(ctx context.Context, authorID string) ([]*domain.Bid, error) {
	return s.Repo.GetBidsByAuthorID(ctx, authorID)
}

func (s *BidService) GetBidStatus(ctx context.Context, bidID, username string) (string, error) {
	// Получаем информацию о ставке
	bid, err := s.Repo.GetBidByID(ctx, bidID)
	if err != nil {
		return "", fmt.Errorf("failed to get bid: %w", err)
	}

	// Проверяем, имеет ли пользователь доступ к этой ставке
	hasAccess, err := s.Repo.UserHasAccessToBid(ctx, username, bidID)
	if err != nil {
		return "", fmt.Errorf("failed to check user access: %w", err)
	}
	if !hasAccess {
		return "", fmt.Errorf("user is not authorized to view this bid status")
	}

	return bid.Status, nil
}

func (s *BidService) EditBid(ctx context.Context, bidID, username, name, description string) error {
	// Проверяем, существует ли заявка
	fmt.Println("bidID:", bidID)
	bid, err := s.Repo.GetBidByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("не удалось получить заявку: %w", err)
	}

	// Проверяем, имеет ли пользователь право редактировать эту заявку
	hasAccess, err := s.Repo.UserHasAccessToBid(ctx, username, bidID)
	if err != nil {
		return fmt.Errorf("не удалось проверить права доступа пользователя: %w", err)
	}
	if !hasAccess {
		return fmt.Errorf("пользователь не имеет прав на редактирование этой заявки")
	}

	// Обновляем данные заявки
	bid.Name = name
	bid.Description = description

	// Сохраняем обновленную заявку в репозитории
	err = s.Repo.UpdateBid(ctx, bid)
	if err != nil {
		return fmt.Errorf("не удалось обновить заявку: %w", err)
	}

	return nil
}