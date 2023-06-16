package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/requests"
)

type GitHubAPIClient struct {
	httpClient    http.Client
	connectionURL string
	authToken     string
}

func NewGitHubAPIClient(connectionURL, authToken string) GitHubAPIClient {
	return GitHubAPIClient{
		httpClient:    *http.DefaultClient,
		connectionURL: connectionURL,
		authToken:     authToken,
	}
}

func (gh GitHubAPIClient) GetPinnedRepositories(username string) ([]PinnedItem, error) {
	requestBody := fmt.Sprintf(gqlQueryString, username)

	responseBytes, err := requests.Query(gh.httpClient, gh.authToken, http.MethodPost, gh.connectionURL, []byte(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}

	var gitHubResponse *GitHubResponse

	err = json.Unmarshal(responseBytes, &gitHubResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %s", err)
	}

	return gitHubResponse.Data.User.PinnedItems.Nodes, nil
}
