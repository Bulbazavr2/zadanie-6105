package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tender_srevice/internal/domain"
)

type TenderRepository interface {
	InsertTender(ctx context.Context, item *domain.Tender) error
	// Добавьте сюда другие методы, которые могут понадобиться в будущем
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (r *PostgresRepository) GetUserIDByUsername(ctx context.Context, username string) (string, error) {
	query := `SELECT id FROM employee WHERE username = $1`
	var userID string
	err := r.DB.QueryRowContext(ctx, query, username).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user ID: %w", err)
	}
	return userID, nil
}


func (r *PostgresRepository) InsertTender(ctx context.Context, item *domain.Tender) error {
	query := `INSERT INTO tenders (name, description, status, service_type, organization_id, creator_username) 
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id`
   

	err := r.DB.QueryRowContext(ctx, query, 
		item.Name, 
		item.Description, 
		item.Status, 
		item.ServiceType, 
		item.OrganizationID, 
		item.CreatorUsername).Scan(&item.ID)
	if err != nil {
		log.Printf("Error inserting tender: %v", err)
		return fmt.Errorf("failed to insert tender: %w", err)
	}
	return nil
}

// Добавить сортировку по алфавиту
func (r *PostgresRepository) GetAllTenders(ctx context.Context) ([]*domain.Tender, error) {
	query := `SELECT id, name, description, status, service_type, organization_id, creator_username FROM tenders`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tenders: %w", err)
	}
	defer rows.Close()

	var tenders []*domain.Tender
	for rows.Next() {
		var t domain.Tender
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.OrganizationID, &t.CreatorUsername); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func (r *PostgresRepository) GetTenderByID(ctx context.Context, tenderID string) (*domain.Tender, error) {
	query := `SELECT id, name, description, status, service_type, organization_id, creator_username
              FROM tenders 
              WHERE id = $1`  
	var t domain.Tender
	err := r.DB.QueryRowContext(ctx, query, tenderID).Scan(
		&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.OrganizationID, &t.CreatorUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tender not found")
		}
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}
	return &t, nil
}

func (r *PostgresRepository) UpdateTenderStatus(ctx context.Context, tender *domain.Tender) error {
	query := `UPDATE tenders 
              SET status = $2
              WHERE id = $1`  
	_, err := r.DB.ExecContext(ctx, query, tender.ID, tender.Status)
    
	if err != nil {
		return fmt.Errorf("failed to update tender status: %w", err)
	}
    
	return nil
}

func (r *PostgresRepository) UpdateTender(ctx context.Context, tender *domain.Tender) error {
	query := `UPDATE tenders 
			  SET name = $2, description = $3, status = $4, service_type = $5, organization_id = $6, creator_username = $7
			  WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query,
		tender.ID, tender.Name, tender.Description, tender.Status,
		tender.ServiceType, tender.OrganizationID, tender.CreatorUsername)
	if err != nil {
		return fmt.Errorf("failed to update tender: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetTendersByUsername(ctx context.Context, username string) ([]*domain.Tender, error) {
	query := `SELECT id, name, description, status, service_type, organization_id, creator_username 
              FROM tenders 
              WHERE creator_username = $1
              ORDER BY name ASC`
	rows, err := r.DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query tenders: %w", err)
	}
	defer rows.Close()

	var tenders []*domain.Tender
	for rows.Next() {
		var t domain.Tender
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.OrganizationID, &t.CreatorUsername); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func (r *PostgresRepository) GetTenderStatus(ctx context.Context, tenderID string) (string, error) {
	query := `SELECT status FROM tenders WHERE id = $1`
	var status string
	err := r.DB.QueryRowContext(ctx, query, tenderID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("tender not found")
		}
		return "", fmt.Errorf("failed to get tender status: %w", err)
	}
	return status, nil
}

//Общий
func (r *PostgresRepository) IsUserResponsibleForOrganization(ctx context.Context, username, organizationID string) (bool, error) {
	// Сначала получаем ID пользователя по его имени
	var userID string
	userID, err := r.GetUserIDByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("ошибка при получении ID пользователя: %w", err)
	}

	// Теперь проверяем, является ли пользователь ответственным за организацию
	query := `SELECT EXISTS(SELECT 1 FROM organization_responsible 
              WHERE user_id = $1 AND organization_id = $2)`
	var exists bool
	err = r.DB.QueryRowContext(ctx, query, userID, organizationID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке ответственности пользователя: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) GetOrganizationIDByTenderID(ctx context.Context, tenderID string) (string, error) {
	query := `SELECT organization_id FROM tenders WHERE id = $1`
	var organizationID string
	err := r.DB.QueryRowContext(ctx, query, tenderID).Scan(&organizationID)
	if err != nil {
		return "", fmt.Errorf("failed to get organization ID: %w", err)
	}
	return organizationID, nil
}

func (r *PostgresRepository) RollbackTender(ctx context.Context, tenderID string, version int) (*domain.Tender, error) {
	query := `
		WITH previous_version AS (
			SELECT name, description, status, service_type, version
			FROM tender_versions
			WHERE tender_id = $1 AND version = $2
		)
		UPDATE tenders t
		SET name = pv.name,
			description = pv.description,
			status = pv.status,
			service_type = pv.service_type,
			version = pv.version
		FROM previous_version pv
		WHERE t.id = $1
		RETURNING t.id, t.name, t.description, t.status, t.service_type, t.organization_id, t.creator_username, t.version
	`

	var t domain.Tender
	err := r.DB.QueryRowContext(ctx, query, tenderID, version).Scan(
		&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.OrganizationID, &t.CreatorUsername, &t.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version not found")
		}
		return nil, fmt.Errorf("failed to rollback tender: %w", err)
	}

	return &t, nil
}

func (r *PostgresRepository) InsertBid(ctx context.Context, bid *domain.Bid) error {
	query := `INSERT INTO bid (name, description, status, tender_id, author_type, author_id)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id, version, created_at`

	err := r.DB.QueryRowContext(ctx, query,
		bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID).
		Scan(&bid.ID, &bid.Version, &bid.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert bid: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetBidsByTenderID(ctx context.Context, tenderID string) ([]*domain.Bid, error) {
	query := `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
			  FROM bid
			  WHERE tender_id = $1
			  ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, tenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bids: %w", err)
	}
	defer rows.Close()

	var bids []*domain.Bid
	for rows.Next() {
		var b domain.Bid
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderID, &b.AuthorType, &b.AuthorID, &b.Version, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, &b)
	}
	return bids, nil
}

// Добавим новый метод для проверки принадлежности пользователя к организации
func (r *PostgresRepository) IsUserInTenderOrganization(ctx context.Context, username, tenderID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tenders t
			JOIN organization_responsible org_resp ON t.organization_id = org_resp.organization_id
			JOIN employee e ON org_resp.user_id = e.id
			WHERE t.id = $1 AND e.username = $2
		)
	`
	var exists bool
	err := r.DB.QueryRowContext(ctx, query, tenderID, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user organization: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) UpdateBidStatus(ctx context.Context, bidID, newStatus string) error {
	query := `UPDATE bid SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, newStatus, bidID)
	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetBidsByAuthorID(ctx context.Context, authorID string) ([]*domain.Bid, error) {
	query := `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
			  FROM bid
			  WHERE author_id = $1
			  ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bids: %w", err)
	}
	defer rows.Close()

	var bids []*domain.Bid
	for rows.Next() {
		var b domain.Bid
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderID, &b.AuthorType, &b.AuthorID, &b.Version, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, &b)
	}
	return bids, nil
}

// GetBidByID возвращает ставку по её ID
func (r *PostgresRepository) GetBidByID(ctx context.Context, bidID string) (*domain.Bid, error) {
	query := `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
			  FROM bid
			  WHERE id = $1`

	var b domain.Bid
	err := r.DB.QueryRowContext(ctx, query, bidID).Scan(
		&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderID, &b.AuthorType, &b.AuthorID, &b.Version, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bid not found")
		}
		return nil, fmt.Errorf("failed to get bid: %w", err)
	}
	return &b, nil
}

func (r *PostgresRepository) IsUserResponsibleForTender(ctx context.Context, username, tenderID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tenders t
			JOIN organization_responsible org_resp ON t.organization_id = org_resp.organization_id
			JOIN employee e ON org_resp.user_id = e.id
			WHERE t.id = $1 AND e.username = $2
		)
	`
	var isResponsible bool
	err := r.DB.QueryRowContext(ctx, query, tenderID, username).Scan(&isResponsible)
	if err != nil {
		return false, fmt.Errorf("failed to check user responsibility: %w", err)
	}
	
	return isResponsible, nil
}

func (r *PostgresRepository) UserHasAccessToBid(ctx context.Context, username, bidID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM bid b
			JOIN tenders t ON b.tender_id = t.id
			JOIN organization_responsible org_resp ON t.organization_id = org_resp.organization_id
			JOIN employee e ON org_resp.user_id = e.id
			WHERE b.id = $1 AND e.username = $2
		)
	`
	var hasAccess bool
	err := r.DB.QueryRowContext(ctx, query, bidID, username).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("failed to check user access to bid: %w", err)
	}
	return hasAccess, nil
}

func (r *PostgresRepository) UpdateBid(ctx context.Context, bid *domain.Bid) error {
	query := `
		UPDATE bid 
		SET name = $2, 
			description = $3, 
			status = $4
		WHERE id = $1
		RETURNING version
	`
	
	err := r.DB.QueryRowContext(ctx, query, 
		bid.ID, 
		bid.Name, 
		bid.Description, 
		bid.Status).Scan(&bid.Version)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("заявка с ID %s не найдена", bid.ID)
		}
		return fmt.Errorf("не удалось обновить заявку: %w", err)
	}

	return nil
}



