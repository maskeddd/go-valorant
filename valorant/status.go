package valorant

import (
	"context"
	"net/http"
	"time"
)

type StatusService service

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgess MaintenanceStatus = "in_progress"
	MaintenanceStatusComplete  MaintenanceStatus = "complete"
)

type IncidentSeverity string

const (
	IncidentSeverityInfo     IncidentSeverity = "info"
	IncidentSeverityWarning  IncidentSeverity = "warning"
	IncidentSeverityCritical IncidentSeverity = "critical"
)

type PublishLocation string

const (
	PublishLocationRiotClient PublishLocation = "riotclient"
	PublishLocationRiotStatus PublishLocation = "riotstatus"
	PublishLocationGame       PublishLocation = "game"
)

type PlatformData struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Locales      []string `json:"locales"`
	Maintenances []Status `json:"maintenances"`
	Incidents    []Status `json:"incidents"`
}

type Status struct {
	ID                int                `json:"id"`
	MaintenanceStatus *MaintenanceStatus `json:"maintenance_status"`
	IncidentSeverity  *IncidentSeverity  `json:"incident_severity"`
	Titles            []StatusContent    `json:"titles"`
	Updates           []Update           `json:"updates"`
	CreatedAt         time.Time          `json:"created_at"`
	ArchiveAt         *time.Time         `json:"archive_at"`
	UpdatedAt         *time.Time         `json:"updated_at"`
	Platforms         []string           `json:"platforms"`
}

type StatusContent struct {
	Locale  string `json:"locale"`
	Content string `json:"content"`
}

type Update struct {
	ID               int               `json:"id"`
	Author           string            `json:"author"`
	Publish          bool              `json:"publish"`
	PublishLocations []PublishLocation `json:"publish_locations"`
	Translations     []StatusContent   `json:"translations"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// ListPlatformData lists the Valorant status for the region's platform.
//
// Valorant API docs: https://developer.riotgames.com/apis#val-status-v1/GET_getPlatformData
func (s *StatusService) ListPlatformData(ctx context.Context) (*PlatformData, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "status/v1/platform-data", nil)
	if err != nil {
		return nil, nil, err
	}

	var platformData *PlatformData
	resp, err := s.client.Do(ctx, req, &platformData)
	if err != nil {
		return nil, resp, err
	}

	return platformData, resp, nil
}
