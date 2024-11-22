package routes

import (
	"TartaLette/api/handlers"
	"github.com/labstack/echo/v4"
)

/*
Fetch -> Get from external source [GitHub]
GET -> Get from internal source [DB, local data]
*/

func RegisterRoutes(e *echo.Echo) {
	// GET general
	e.GET("/health", handlers.HealthCheck)

	// GET GroupState
	e.GET("/state/:hostname", handlers.GetRunnersForHostname)
	e.GET("/state/hostname", handlers.GetAvailableHostnames)
	e.GET("/state/:hostname/:runnerName", handlers.GetRunnerInfo)

	// GET GroupGitHub
	e.GET("/gh/runner", handlers.FetchRunnerInfoByName)           // GH & DB
	e.GET("/gh/fetchState", handlers.FetchCurrentStateFromGithub) // GH & DB
	// POST GroupGitHub
	e.POST("/gh/webhook", handlers.WebhookStateUpdateFromGitHub) // GH & DB

	// GET GroupRunner
	e.GET("/runner/register", handlers.FetchNewRegistrationToken) // GH & DB

	// Get Static
	e.File("/favicon.ico", "static/images/favicon.ico")
}
