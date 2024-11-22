package utils

import (
	"TartaLette/api/models"
	"TartaLette/gh"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func ReadConfigs(client *gh.Client, filesToRead []string) ([]models.TartState, error) {
	// Re-read config files in the specified directory
	// This function should be called upon start, and then periodically
	// It should also be called when a new config is added, by a webhook from GH itself.

	var tartConfigs []models.TartState

	for _, file := range filesToRead {
		logrus.Info(fmt.Sprintf("Reading file [%v]", file))

		var tartConfig models.TartState
		content, err := client.FetchFileContent(file)

		if err != nil {
			logrus.Error(fmt.Sprintf("Couldn't fetch file [%v] content -> \n%v", file, err))
		}
		base64content, err := base64.StdEncoding.DecodeString(content)

		if err != nil {
			logrus.Error(fmt.Sprintf("Couldn't decode file [%v] content -> \n%v", file, err))
		}

		err = yaml.Unmarshal(base64content, &tartConfig)
		if err != nil {
			logrus.Error(fmt.Sprintf("Couldn't unmarshal file [%v] content -> \n%v", file, err))
			logrus.Error("Are you sure structure is correct?")
		}

		logrus.Debug(fmt.Sprintf("File [%+v] content -> \n%v", file, tartConfig))

		// state/pa4-r102-srv-macrunner02 => pa4-r102-srv-macrunner02
		tartConfig.HostName = strings.Split(file, "/")[1]
		tartConfigs = append(tartConfigs, tartConfig)
	}

	return tartConfigs, nil

}

func FetchState(client *gh.Client) ([]models.Hostname, []models.Runner) {

	var hostnames []models.Hostname
	var runners []models.Runner
	var files []string
	var stateFolder string

	stateFolder = os.Getenv("STATE_DIRECTORY")

	if stateFolder == "" {
		stateFolder = "state"
	}
	files, err := client.ListDirectory(stateFolder)
	if err != nil {

	}

	data, err := ReadConfigs(client, files)

	if err != nil {
		logrus.Fatal(fmt.Sprintf("Couldn't read configs -> \n%v", err))
	}

	for _, tartState := range data {
		hostname := models.Hostname{HostName: tartState.HostName}
		hostnames = append(hostnames, hostname)
		for _, runner := range tartState.Runners {
			runner.HostName = tartState.HostName
			// Since runners can have the same name, we need to make sure they are unique, at least on host level
			runner.Name = runner.Name + ":" + tartState.HostName
			runners = append(runners, runner)
		}
	}
	return hostnames, runners
}
