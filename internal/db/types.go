package db

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
)

type DBConnector interface {
	CreateCommitEntry(timeline models.Timelines) error
	GetCommitsAll(offset int, limit int) ([]models.Timelines, int64, error)
	GetCommitsByService(service models.Services, offset int, limit int) ([]models.Timelines, int64, error)
	GetCommitByRef(ref string) (models.Timelines, int64, error)

	CreateServiceTableEntry(name string, s config.Service) (models.Services, error)
	GetServicesAll(offset int, limit int) ([]models.ExpandedServices, int64, error)
	GetLatest(service models.Services) (models.ExpandedServices, error, error)
	GetServiceByName(name string) (models.Services, error)
	GetServiceByGHRepo(repo string) (models.Services, error)

	GetTimelinesAll(offset int, limit int) ([]models.Timelines, int64, error)
	GetTimelinesByService(service models.Services, offset int, limit int) ([]models.Timelines, int64, error)
	GetTimelineByRef(ref string) (models.Timelines, int64, error)

	GetDeploysAll(offset int, limit int) ([]models.Timelines, int64, error)
	GetDeploysByService(service models.Services, offset int, limit int) ([]models.Timelines, int64, error)
	GetDeployByRef(ref string) (models.Timelines, int64, error)
}
