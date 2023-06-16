package model

type ProjectItem struct {
	ID              string `json:"_id,omitempty"`
	Articlelink     string `json:"articlelink"`
	Codelink        string `json:"codelink"`
	Demolink        string `json:"demolink"`
	Details         string `json:"details"`
	DisplayPriority string `json:"displayPriority"`
	Domain          string `json:"domain"`
	Name            string `json:"name"`
	Stack           string `json:"stack"`
}
