package handlers

import (
	"server/config"

	"github.com/labstack/echo/v4"
)

// RunMigration godoc
// @Summary Execute a migration
// @Description Execute a specific database migration identified by its ID
// @Tags Migrations
// @Accept json
// @Produce json
// @Param migration_id body models.RunMigrationRequest true "Migration ID to run"
// @Success 200 {object} map[string]string "Migration ran successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Migration not found"
// @Failure 500 {object} map[string]string "Error running migration"
// @Router /api/v1/admin/run-migrations [post]
func RunMigration(c echo.Context) error {
	return config.RunMigration(c)
}

// GetMigration godoc
// @Summary Retrieve all available migrations
// @Description Get a list of all available migrations with their titles and descriptions
// @Tags Migrations
// @Accept json
// @Produce json
// @Success 200 {array} models.GetMigrationListRequest "List of migrations"
// @Failure 500 {object} map[string]string "Failed to load migrations"
// @Router /api/v1/admin/get-migrations [get]
func GetMigration(c echo.Context) error {
	return config.GetMigration(c)
}
