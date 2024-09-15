package domain

const (
	TenderStatusCreated   = "CREATED"
	TenderStatusPublished = "PUBLISHED"
	TenderStatusClosed    = "CLOSED"
)

type Tender struct {
	ID               string  `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ServiceType      string `json:"serviceType"`
	Version          int    `json:"version"`
	Status           string `json:"status"`
	OrganizationID   string  `json:"organizationId"`
	CreatorUsername  string `json:"creatorUsername"`
}