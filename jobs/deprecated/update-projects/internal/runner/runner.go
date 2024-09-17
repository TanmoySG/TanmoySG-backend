package runner

import (
	"encoding/json"
	"fmt"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/common"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/github"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/transform"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/wdb"

	log "github.com/sirupsen/logrus"
)

type client struct {
	wdbAdapter   wdb.WdbAdapter
	githubClient github.GitHubAPIClient
}

type configurations struct {
	wdbDatabase   string
	wdbCollection string
	username      string
}

func NewClient(wdbAdapter wdb.WdbAdapter, githubClient github.GitHubAPIClient) client {
	return client{
		wdbAdapter:   wdbAdapter,
		githubClient: githubClient,
	}
}

func NewConfig(database, collection, username string) configurations {
	return configurations{
		wdbDatabase:   database,
		wdbCollection: collection,
		username:      username,
	}
}

func (c client) Run(configurations configurations) error {
	pinnedRepositories, err := c.githubClient.GetPinnedRepositories(configurations.username)
	if err != nil {
		return fmt.Errorf("error getting pinned repositories: %s", err)
	}
	log.Infof("fetched pinned repositories for user [%s] from gitHub.  total repositories fetched [%d]", configurations.username, len(pinnedRepositories))

	existingProjectsMap, err := c.wdbAdapter.GetData(configurations.wdbDatabase, configurations.wdbCollection)
	if err != nil {
		return fmt.Errorf("error getting existing projects map: %s", err)
	}

	existingProjects, err := common.GetProjectsList(existingProjectsMap.Data)
	if err != nil {
		return fmt.Errorf("error getting existing projects list: %s", err)
	}

	log.Infof("fetched existing projects. total projects fetched [%d]", len(existingProjects))

	for sequence, item := range pinnedRepositories {
		actualSequence := len(pinnedRepositories) - sequence - 1
		transformedProject := transform.Do(actualSequence, item, existingProjects)
		log.Infof("transformed repo details to project [%s / %s], sequence [%d]", transformedProject.Name, transformedProject.Codelink, actualSequence)

		transformedProjectByte, err := json.Marshal(transformedProject)
		if err != nil {
			return fmt.Errorf("error marshalling transformedProject to bytes: %s", err)
		}

		var transformedProjectMap map[string]interface{}
		err = json.Unmarshal(transformedProjectByte, &transformedProjectMap)
		if err != nil {
			return fmt.Errorf("error unmarshalling transformedProject bytes to map: %s", err)
		}

		err = c.wdbAdapter.AddData(configurations.wdbDatabase, configurations.wdbCollection, transformedProjectMap)
		if err != nil {
			return fmt.Errorf("error adding data to wdb-collection: %s", err)
		}
		log.Infof("inserted pinned repo as project [%s]", transformedProject.Name)
	}

	for _, existingProject := range existingProjects {
		err = c.wdbAdapter.DeleteData(configurations.wdbDatabase, configurations.wdbCollection, "_id", existingProject.ID)
		if err != nil {
			return fmt.Errorf("error deleting data from wdb-collection: %s", err)
		}
		log.Infof("deleted project [%s]", existingProject.Name)
	}

	return nil
}
