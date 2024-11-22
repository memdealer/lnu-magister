// Package Gh: Github API Client wrapper
// Tests: Gh/client_test.go (TODO)

package gh

import (
	"context"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
	"github.com/sirupsen/logrus"
	"net/http"
)

/*
Abstract:
Get registration token to register runners; [DONE]
Get the file content, pass it return it as a byte array. [DONE]
List files in a directory. [DONE]
Get runner group ID and information. [DONE]
Get runners in runner group. [DONE]
Get the runner, check if it exists and its status. [DONE]
*/

type Client struct {
	GhClient       *github.Client
	Organization   string
	RepositoryName string
}

func NewClient(appID int64, installationID int64, organization string, repositoryName string, privateKeyPath string) Client {

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, privateKeyPath)

	if err != nil {
		logrus.Fatal("Error: %v\n", err)
	}
	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})

	return Client{GhClient: client, Organization: organization, RepositoryName: repositoryName}
}

func (c *Client) FetchFileContent(filePath string) (string, error) {
	file, _, _, err := c.GhClient.Repositories.GetContents(context.Background(), c.Organization, c.RepositoryName, filePath, nil)
	if err != nil {
		logrus.Warn("Errored during GetFileContent call, or file does not exist:\n", err)
		return "", err
	}
	return *file.Content, nil
}

func (c *Client) ListDirectory(dirPath string) ([]string, error) {
	var files []string

	_, dirs, _, err := c.GhClient.Repositories.GetContents(context.Background(), c.Organization, c.RepositoryName, dirPath, nil)
	if err != nil {
		logrus.Warn("Errored during ListDirectory call, or directory does not exist:\n", err)
		return nil, err
	}
	for _, file := range dirs {
		files = append(files, *file.Path)
	}

	return files, nil
}

func (c *Client) GetRegistrationToken() (string, error) {
	token, _, err := c.GhClient.Actions.CreateOrganizationRegistrationToken(context.Background(), c.Organization)
	if err != nil {
		logrus.Warn("Errored during GeneratePatToken call:\n", err)
		return "", err
	}
	return *token.Token, nil
}

func (c *Client) FindRunnerGroupByName(runnerGroupName string) (github.RunnerGroup, error) {
	runnerGroups, _, err := c.GhClient.Actions.ListOrganizationRunnerGroups(context.Background(), c.Organization, nil)
	if err != nil {
		logrus.Warn("Errored during FindRunnerGroupID call:\n", err)
		return github.RunnerGroup{}, err
	}
	for _, runnerGroup := range runnerGroups.RunnerGroups {
		if *runnerGroup.Name == runnerGroupName {
			return *runnerGroup, nil
		}
	}
	return github.RunnerGroup{}, nil
}

func (c *Client) ListRunnersInRunnerGroup(runnerGroupId int64) ([]*github.Runner, error) {
	runners, _, err := c.GhClient.Actions.ListRunnerGroupRunners(context.Background(), c.Organization, runnerGroupId, nil)
	if err != nil {
		logrus.Warn("Errored during ListRunnersInRunnerGroup call:\n", err)
		return nil, err
	}

	return runners.Runners, nil
}

func (c *Client) FindRunnerByNameWithPager(runnerName string) (*github.Runner, error) {
	opt := &github.ListOptions{PerPage: 100}
	for {
		runners, resp, err := c.GhClient.Actions.ListOrganizationRunners(context.Background(), c.Organization, opt)
		if err != nil {
			logrus.Warn("Errored during FindRunnerByNameWithPager call:\n", err)
			return nil, err
		}
		for _, runner := range runners.Runners {
			if *runner.Name == runnerName {
				return runner, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil, nil
}
