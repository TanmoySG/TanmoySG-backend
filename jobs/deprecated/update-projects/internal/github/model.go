package github

type GitHubResponse struct {
	Data struct {
		User struct {
			PinnedItems struct {
				Nodes []PinnedItem `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"user"`
	} `json:"data"`
}

type PinnedItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Languages   struct {
		Edges []struct {
			Node struct {
				Name string `json:"name"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"languages"`
	RepositoryTopics struct {
		Edges []struct {
			Node struct {
				Topic struct {
					Name string `json:"name"`
				} `json:"topic"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"repositoryTopics"`
	HomepageURL string `json:"homepageUrl"`
}
