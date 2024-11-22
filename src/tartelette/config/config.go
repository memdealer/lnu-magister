// Package Config File: config.go
// Configuration manager.
// This file is responsible for loading the configuration from the environment variables
// Tests: Config/config_test.go (TODO)

package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
)

var AppConfig Configuration

type Configuration struct {
	GithubAppId             int64
	GitHubInstallationId    int64
	GithubAppPrivateKeyPath string
	StateDirectory          string
	RepositoryName          string
	Organization            string
	RunnerGroupName         string
	ApiKey                  string
	GitHubWebHookSecret     string
}

func getEnvConfigInt(envName string) int64 {
	//To convert string to int,
	//The os.getenv() returns only string, we need to convert it to int

	envVar := os.Getenv(envName)

	if envVar == "" {
		envVar = "0"
	}

	value, err := strconv.ParseInt(envVar, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func init() {
	// Construct, validate and return a configuration if it is not empty

	requiredVars := "GITHUB_APP_ID, GITHUB_INSTALLATION_ID, " +
		"GITHUB_APP_PRIVATE_KEY_PATH, STATE_DIRECTORY," +
		"REPOSITORY_NAME, GITHUB_ORGANIZATION, GITHUB_RUNNER_GROUP_NAME, ApiKey, GITHUB_WEBHOOK_SECRET"
	logrus.Println("Loading configuration...")

	err := godotenv.Load("test.env")
	if err != nil {
		logrus.Warn("No .env file found, loading from environment variables")
	}

	AppConfig = Configuration{
		GithubAppId:             getEnvConfigInt("GITHUB_APP_ID"),
		GitHubInstallationId:    getEnvConfigInt("GITHUB_INSTALLATION_ID"),
		GithubAppPrivateKeyPath: os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"),
		StateDirectory:          os.Getenv("STATE_DIRECTORY"),
		RepositoryName:          os.Getenv("REPOSITORY_NAME"),
		Organization:            os.Getenv("GITHUB_ORGANIZATION"),
		RunnerGroupName:         os.Getenv("GITHUB_RUNNER_GROUP_NAME"),
		ApiKey:                  os.Getenv("API_KEY"),
		GitHubWebHookSecret:     os.Getenv("GITHUB_WEBHOOK_SECRET"),
	}

	// I did not use Reflect, or anything fancy to validate the configuration
	// As using Reflect is like using a shotgun whilst having a bathtub. Kinda works, but not really.
	if AppConfig.GithubAppId == 0 || AppConfig.GitHubInstallationId == 0 {
		logrus.Fatal(fmt.Sprintf(`GITHUB_APP_ID or GITHUB_INSTALLATION_ID is not set, required variables:
[%s]`, requiredVars))
	}

	if AppConfig.GithubAppPrivateKeyPath == "" {
		logrus.Fatal(fmt.Sprintf(`GITHUB_APP_PRIVATE_KEY_PATH is not set, required variables:
[%s]`, requiredVars))
	}

	if AppConfig.StateDirectory == "" {
		logrus.Fatal(fmt.Sprintf(`STATE_DIRECTORY is not set, required variables:
[%s]`, requiredVars))
	}

	if AppConfig.RepositoryName == "" {
		logrus.Fatal(fmt.Sprintf(`REPOSITORY_NAME is not set, required variables:
[%s]`, requiredVars))
	}

	if AppConfig.Organization == "" {
		logrus.Fatal(fmt.Sprintf(`ORGANIZATION is not set, required variables:
[%s]`, requiredVars))

	}
	if AppConfig.RunnerGroupName == "" {
		logrus.Fatal(fmt.Sprintf(`RunnerGroupName is not set, required variables:
[%s]`, requiredVars))

	}

	if AppConfig.ApiKey == "" {
		logrus.Fatal(fmt.Sprintf(`API_KEY is not set, required variables: 
[%s]`, requiredVars))
	}

	if AppConfig.GitHubWebHookSecret == "" {
		logrus.Fatal(fmt.Sprintf(`GITHUB_WEBHOOK_SECRET is not set, required variables:
[%s]`, requiredVars))
	}

	logrus.Info("Configuration is valid")
	logrus.Info(fmt.Sprintf("Operating on: [%v]\n Repository: [%v]\n RunnerGroup: [%v]",
		AppConfig.Organization, AppConfig.RepositoryName, AppConfig.RunnerGroupName))
	logrus.Info("Configuration Valid - [OK]")

}
