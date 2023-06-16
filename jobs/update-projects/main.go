package main

import (
	"os"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/github"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/runner"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/wdb"
	log "github.com/sirupsen/logrus"
)

func main() {
	baseURL, cluster, token := os.Getenv("WDB_RETRO_BASE_URL"), os.Getenv("WDB_RETRO_CLUSTER"), os.Getenv("WDB_RETRO_TOKEN")
	wdbAdapter := wdb.NewClient(baseURL, cluster, token)

	authToken, queryURL := os.Getenv("GH_AUTH_TOKEN"), os.Getenv("GH_QUERY_URL")
	ghClient := github.NewGitHubAPIClient(queryURL, authToken)

	usernameToQuery, wdbDatabase, wdbCollection := os.Getenv("GH_USERNAME"), os.Getenv("WDB_DATABASE"), os.Getenv("WDB_COLLECTION")
	conf := runner.NewConfig(wdbDatabase, wdbCollection, usernameToQuery)

	rc := runner.NewClient(wdbAdapter, ghClient)
	
	err := rc.Run(conf)
	if err != nil {
		log.Fatalf(" Aborting operation, error while updating projects: %s", err)
	}
}
