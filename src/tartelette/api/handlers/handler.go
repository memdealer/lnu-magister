package handlers

import (
	"TartaLette/api/models"
	"TartaLette/gh"
	"TartaLette/utils"
	"github.com/google/go-github/v53/github"
	"github.com/hashicorp/go-memdb"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const standartResponse = "SOSI HOOI, BYDLO."

func HealthCheck(e echo.Context) error {
	// Basic health check that server is UP and running, no extra magic aroud here.
	return e.JSON(http.StatusOK, "OK")
}

func FetchNewRegistrationToken(c echo.Context) error {
	// See the name.
	var result models.Token
	ghClient := c.Get("ghClient").(*gh.Client)

	tartRunnerHeader := c.Request().Header.Get("X-Tart-Runner-Name")
	logrus.Warn("Got a request for token for runner: ", tartRunnerHeader)
	// For tracking purpuses, we need to know which runner is asking for a token
	if tartRunnerHeader == "" {
		return c.JSON(http.StatusBadRequest, "Missing X-Tart-Runner-Name header")
	}

	data, err := ghClient.GetRegistrationToken()

	result.Token = data

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}

func FetchCurrentStateFromGithub(c echo.Context) error {
	// Might be useful to re-reading state once in a while, or forcefully doing so.

	ghClient := c.Get("ghClient").(*gh.Client)
	dbConnection := c.Get("dbConnection").(*memdb.MemDB)

	hostnames, runners := utils.FetchState(ghClient)

	if len(hostnames) == 0 {
		return c.JSON(http.StatusInternalServerError, "No hostnames found")
	}
	if len(runners) == 0 {
		return c.JSON(http.StatusInternalServerError, "No runners found")
	}

	err := utils.FillDbWithStateValuesFromGithub(dbConnection, hostnames, runners)
	if err != nil {
		logrus.Error("Couldn't update database")
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "OK")

}

func FetchRunnerInfoByName(c echo.Context) error {
	// Fetch the status of a runner by name.
	ghClient := c.Get("ghClient").(*gh.Client)
	runnerName := c.QueryParam("runnerName")

	if runnerName == "" {
		return c.JSON(http.StatusBadRequest, "Missing runner name")
	}

	runner, err := ghClient.FindRunnerByNameWithPager(runnerName)

	// Check if the runner was found
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if runner == nil {
		return c.JSON(http.StatusNotFound, "Runner not found")
	}

	return c.JSON(http.StatusOK, runner)
}

func GetRunnersForHostname(c echo.Context) error {
	// /state/:hostname -> RunnerList

	var runners []models.Runner

	hostname := c.Param("hostname")
	dbConnection := c.Get("dbConnection").(*memdb.MemDB)
	txn := dbConnection.Txn(false)
	defer txn.Abort()

	if hostname == "" {
		return c.JSON(http.StatusBadRequest, "Missing hostname")
	}

	it, err := txn.Get("runner", "hostName", hostname)

	if err != nil {
		logrus.Error("Couldn't read database")
		return c.JSON(http.StatusInternalServerError, err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		runner := obj.(models.Runner)
		runners = append(runners, runner)

	}

	if len(runners) == 0 {
		return c.JSON(http.StatusNotFound, "No runners found for this hostname")

	}
	return c.JSON(http.StatusOK, runners)

}

func GetRunnerInfo(c echo.Context) error {
	// /state/:hostname/:runnerName -> RunnerInfo

	var runner models.Runner

	hostname := c.Param("hostname")
	runnerName := c.Param("runnerName")

	dbConnection := c.Get("dbConnection").(*memdb.MemDB)
	txn := dbConnection.Txn(false)
	defer txn.Abort()

	if hostname == "" {
		return c.JSON(http.StatusBadRequest, "Missing hostname")
	}

	if runnerName == "" {
		return c.JSON(http.StatusBadRequest, "Missing runner name")
	}

	it, err := txn.Get("runner", "hostName", hostname)

	if err != nil {
		logrus.Error("Couldn't read database")
		return c.JSON(http.StatusInternalServerError, err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		runner = obj.(models.Runner)
		if runner.Name == runnerName {
			return c.JSON(http.StatusOK, runner)
		}
	}

	return c.JSON(http.StatusNotFound, "Runner not found")

}

func GetAvailableHostnames(c echo.Context) error {
	// /state -> HostnameList

	var hostnames []models.Hostname

	dbConnection := c.Get("dbConnection").(*memdb.MemDB)
	txn := dbConnection.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("host", "id")

	if err != nil {
		logrus.Error("Couldn't read database")
		return c.JSON(http.StatusInternalServerError, err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		hostname := obj.(models.Hostname)
		hostnames = append(hostnames, hostname)
	}

	if len(hostnames) == 0 {
		return c.JSON(http.StatusNotFound, "No hostnames found")
	}

	return c.JSON(http.StatusOK, hostnames)
}

// /state/:hostname/:runnerName -> RunnerInfo

func WebhookStateUpdateFromGitHub(c echo.Context) error {
	// This is the most intense stuff, refactor required,

	payload, err := ioutil.ReadAll(c.Request().Body)
	header := c.Request().Header.Get("X-Hub-Signature")
	// Check if it is a post request
	if c.Request().Method != "POST" {
		return c.JSON(http.StatusMethodNotAllowed, standartResponse)

	}

	if err != nil {
		logrus.Error("Couldn't read body")
		return c.JSON(http.StatusInternalServerError, err)
	}

	if header == "" {
		logrus.Error("Missing X-Hub-Signature header")
		return c.JSON(http.StatusForbidden, standartResponse)
	}

	if !utils.ValidateSignature(header, payload) {
		logrus.Error("Payload validation failed")
		return c.JSON(http.StatusForbidden, standartResponse)
	}

	// If the signature is valid, we can proceed with the update
	logrus.Info("Payload validation successful")

	// Check that "modified" files from github webhook `push` event are in the `state` folder
	// If not, we don't care about the update
	data, err := github.ParseWebHook(github.WebHookType(c.Request()), payload)
	if err != nil {
		logrus.Error("Couldn't parse webhook")
		return c.JSON(http.StatusInternalServerError, err)
	}

	res, err := utils.CheckIfModifiedFilesAreInStateFolder(data)

	if err != nil {
		logrus.Error("Couldn't check modified files")
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !res {
		logrus.Info("State folder not modified, ignoring update")
		return c.JSON(http.StatusOK, "OK")
	}

	// If the files are in the state folder, we can proceed with the update
	ghClient := c.Get("ghClient").(*gh.Client)
	dbConnection := c.Get("dbConnection").(*memdb.MemDB)

	utils.FetchStateAndCommit(dbConnection, ghClient)
	return c.JSON(http.StatusOK, "OK")
}
