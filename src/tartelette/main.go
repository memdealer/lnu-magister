package main

import (
	"TartaLette/api/middlewares"
	"TartaLette/api/routes"
	"TartaLette/config"
	"TartaLette/db"
	"TartaLette/gh"
	"TartaLette/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ValidatorConfig(key string, c echo.Context) (bool, error) {
	return key == config.AppConfig.ApiKey, nil
}

func main() {

	client := gh.NewClient(config.AppConfig.GithubAppId,
		config.AppConfig.GitHubInstallationId,
		config.AppConfig.Organization,
		config.AppConfig.RepositoryName,
		config.AppConfig.GithubAppPrivateKeyPath)

	dbConnection := db.InitDatabase()

	e := echo.New()

	// Custom middlewares
	e.Use(
		middlewares.GitHub(&client),
	)

	e.Use(
		middlewares.DbConnection(dbConnection),
	)

	// Global, general Middlewares
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "${time_rfc3339}, took:${latency_human}, method:${method}, uri:${uri}, ${status} ${error}\n"}))

	// This one to prevent unauthorized access to the API
	e.Use(middleware.KeyAuthWithConfig(
		middleware.KeyAuthConfig{
			KeyLookup: "header:X-TartaLette-Api-Key",
			Validator: ValidatorConfig,
		}))

	// This one to generate unique request ID
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	// Routes
	routes.RegisterRoutes(e)

	//go utils.GithubReaderHeartBeat(dbConnection, &client)
	utils.FetchStateAndCommit(dbConnection, &client)
	e.Logger.Fatal(e.Start(":1323"))
}
