package common

import (
	"encoding/json"
	"fmt"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/model"
)

func GetProjectsList(data map[string]interface{}) ([]model.ProjectItem, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %s", err)
	}

	var projectsMap *map[string]model.ProjectItem
	err = json.Unmarshal(dataBytes, &projectsMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %s", err)
	}

	projectsList := []model.ProjectItem{}
	for _, item := range *projectsMap {
		projectsList = append(projectsList, item)
	}

	return projectsList, nil
}
