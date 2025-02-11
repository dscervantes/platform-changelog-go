package db

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
)

func (conn *DBConnectorImpl) CreateServiceTableEntry(s *models.Services) error {
	results := conn.db.Create(s)

	return evaluateError(results.Error)
}

func (conn *DBConnectorImpl) UpdateServiceTableEntry(name string, s config.Service) (service models.Services, err error) {
	newService := models.Services{Name: name, DisplayName: s.DisplayName, Tenant: s.Tenant}
	results := conn.db.Model(models.Services{}).Where("name = ?", name).Updates(&newService)

	return newService, evaluateError(results.Error)
}

func (conn *DBConnectorImpl) DeleteServiceTableEntry(name string) (models.Services, error) {
	// save the service to delete the timelines
	service, _, _ := conn.GetServiceByName(name)

	results := conn.db.Model(models.Services{}).Where("name = ?", name).Delete(&models.Services{})
	if results.Error != nil {
		return models.Services{}, evaluateError(results.Error)
	}

	// delete the timelines for the service
	err := conn.DeleteTimelinesByService(service)
	if err != nil {
		return models.Services{}, evaluateError(err)
	}

	return service, nil
}

func (conn *DBConnectorImpl) GetServicesAll(offset int, limit int, q structs.Query) ([]structs.ExpandedServicesData, int64, error) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServicesAll)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var services []structs.ExpandedServicesData

	db := conn.db.Model(models.Services{})

	if len(q.Name) > 0 {
		db = db.Where("services.name IN ?", q.Name)
	}
	if len(q.DisplayName) > 0 {
		db = db.Where("services.display_name IN ?", q.DisplayName)
	}
	if len(q.Tenant) > 0 {
		db = db.Where("services.tenant IN ?", q.Tenant)
	}

	// Uses the Services model here to reflect the proper db relation
	db.Model(models.Services{}).Count(&count)

	// TODO: add a sort_by field to the query struct
	result := db.Order("ID desc").Preload("Projects").Limit(limit).Offset(offset).Scan(&services)

	var servicesWithTimelines []structs.ExpandedServicesData
	for i := 0; i < len(services); i++ {
		s, _, _ := conn.GetLatest(services[i])

		servicesWithTimelines = append(servicesWithTimelines, s)
	}

	return servicesWithTimelines, count, evaluateError(result.Error)
}

func (conn *DBConnectorImpl) GetLatest(service structs.ExpandedServicesData) (structs.ExpandedServicesData, error, error) {
	l.Log.Debugf("Query name: %s", service.Name)

	// TODO: Make one query to get the latest commit and deploy for each service
	comResult := conn.db.Model(models.Timelines{}).Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "commit").Order("Timestamp desc").Limit(1).Find(&service.Commit)

	depResult := conn.db.Model(models.Timelines{}).Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "deploy").Order("Timestamp desc").Limit(1).Find(&service.Deploy)

	return service, evaluateError(comResult.Error), evaluateError(depResult.Error)
}

func (conn *DBConnectorImpl) GetServiceNames() ([]string, error) {
	var names []string
	result := conn.db.Model(models.Services{}).Pluck("name", &names)
	return names, evaluateError(result.Error)
}

func (conn *DBConnectorImpl) GetServiceByName(name string) (models.Services, int64, error) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServiceByName)
	defer callDurationTimer.ObserveDuration()

	var service models.Services
	result := conn.db.Model(models.Services{}).Preload("Projects").Where("services.name = ?", name).First(&service)
	return service, result.RowsAffected, evaluateError(result.Error)
}

func (conn *DBConnectorImpl) GetServiceByRepo(repo string) (models.Services, error) {
	var service models.Services
	result := conn.db.Model(models.Services{}).Preload("Projects").Where("repo = ?", repo).First(&service)

	return service, evaluateError(result.Error)
}
