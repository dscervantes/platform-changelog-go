package main

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
)

func seedDB(cfg *config.Config) {
	logging.Log.Info("Seeding DB")

	// dbConnector := db.NewDBConnector(cfg)

	// cleanupServices(cfg, dbConnector) // Remove services that are no longer in the config
	// reconcileServices(cfg, dbConnector)

	logging.Log.Info("DB Seeding Complete")
}

func cleanupServices(cfg *config.Config, conn *db.DBConnectorImpl) {
	names, _ := conn.GetServiceNames()
	for _, name := range names {
		if _, ok := cfg.Services[name]; !ok {
			conn.DeleteServiceTableEntry(name)
			logging.Log.Info("Deleted service: ", name)
		}
	}
}

func reconcileServices(cfg *config.Config, conn *db.DBConnectorImpl) {
	for key, service := range cfg.Services {
		// Validate the tenant field exists in the config
		if !validateTenant(service.Tenant, cfg) {
			logging.Log.Error("Tenant not validated: ", service.Tenant)
			continue
		}

		serviceData, _, err := conn.GetServiceByName(key)
		if err != nil {
			logging.Log.Error("Error getting service: ", err)
		}

		// update the service
		if compareService(serviceData, service) {
			// if the service in the config and db are different
			err := updateService(key, service, conn)
			if err != nil {
				logging.Log.Error("Error updating service: ", err)
				continue
			}
		}
	}
}

func validateTenant(tenant string, cfg *config.Config) bool {
	for _, t := range cfg.Tenants {
		if t.Name == tenant {
			return true
		}
	}
	return false
}

// Compare the services in the DB to the services in the config
// Returns false is they are the same, true if they are different
func compareService(fromDB models.Services, fromCfg config.Service) bool {
	if fromDB.DisplayName != fromCfg.DisplayName ||
		fromDB.Tenant != fromCfg.Tenant {
		// TODO: bulk update services
		// logging.Log.Info("Queued service for update: ", fromCfg.DisplayName)
		return true
	}

	return false
}

func updateService(name string, fromCfg config.Service, conn *db.DBConnectorImpl) error {
	_, err := conn.UpdateServiceTableEntry(name, fromCfg)
	if err != nil {
		logging.Log.Error("Error updating service: ", err)
		return err
	}

	logging.Log.Info("Updated service: ", fromCfg.DisplayName)
	return nil
}
