package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/google/go-github/v53/github"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func ValidateSignature(signature string, payload []byte) bool {
	secret := []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")) // Replace with your GitHub webhook secret

	mac := hmac.New(sha1.New, secret)
	_, _ = mac.Write(payload)
	expectedHash := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedHash))
}

func CheckIfModifiedFilesAreInStateFolder(data any) (bool, error) {

	var stateFolder string

	stateFolder = os.Getenv("STATE_DIRECTORY")

	if stateFolder == "" {
		stateFolder = "state"
	}

	// Assert head_commit, if it contains /state folder
	// If not, we don't care about the update
	switch event := data.(type) {
	// Switch for later use.
	case *github.PushEvent:
		if event.HeadCommit == nil {
			logrus.Error("Head commit is nil")
			return false, errors.New("head commit is nil")
		}
		for _, modifiedFile := range event.HeadCommit.Modified {
			println(modifiedFile)
			if strings.Split(modifiedFile, "/")[0] == stateFolder {
				logrus.Info("State folder modified, proceeding with update")
				return true, nil
			} else {
				logrus.Info("State folder not modified, ignoring update")
				return false, nil
			}
		}
	}
	return false, errors.New("unknown error")
}
