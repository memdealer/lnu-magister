package utils

import (
	"TartaLette/gh"
	"github.com/hashicorp/go-memdb"
	"github.com/sirupsen/logrus"
	"time"
)

func FetchStateAndCommit(dbConnection *memdb.MemDB, client *gh.Client) {

	hostnames, runners := FetchState(client)

	if len(hostnames) == 0 {
		logrus.Error("No hosts found")
	}
	if len(runners) == 0 {
		logrus.Error("No runners found")
	}

	// Update the database
	err := FillDbWithStateValuesFromGithub(dbConnection, hostnames, runners)
	if err != nil {
		logrus.Error("Couldn't update database")

	}
}

func GithubReaderHeartBeat(dbConnection *memdb.MemDB, client *gh.Client) {
	// Initially, I thought it woulbe nicer idea to use github webhooks to update the database
	// But I think it's better to have a cron job that will update the database every minute
	// This way, we can have a better control of the database and the data that is stored in it
	// And we can also have a better control of the data that is sent to the frontend
	// This function will be called upon start, and then periodically
	// It should also be called when a new config is added, by a webhook from GH itself.
	// One of the reasons against webhook -> is that this particular app has access to organization runners RW
	// therefore exposing it to the internet is not a good idea, as we cannot guarantee rogue access to the runners

	// To reconsider.

	firstRun := false

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		if firstRun == false {
			logrus.Info("Reading Github config initially")
			FetchStateAndCommit(dbConnection, client)
			firstRun = true
		}
		select {
		case <-ticker.C:
			logrus.Info("Reading Github config")
			FetchStateAndCommit(dbConnection, client)
		}
	}
}
