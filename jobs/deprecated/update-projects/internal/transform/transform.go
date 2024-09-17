package transform

import (
	"strconv"
	"strings"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/github"
	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// EDIT This to add more domains
var domains = map[string]string{
	"default":           "web application",
	"cryptography":      "cryptography",
	"command line tool": "command line tool",
	"cli":               "command line tool",
}

func Do(sequence int, pinnedItem github.PinnedItem, existingProjects []model.ProjectItem) model.ProjectItem {
	sequenceString := strconv.Itoa(sequence)

	domain, articlelink := getExistingValues(pinnedItem, existingProjects)
	if domain == "" {
		domain = enrichDomain(pinnedItem)
	}

	if articlelink == "" {
		articlelink = pinnedItem.URL
	}

	demoLink := pinnedItem.HomepageURL
	if demoLink == "" {
		demoLink = "NULL"
	}

	stack := enrichStack(pinnedItem)
	caser := cases.Title(language.English)
	projectItem := model.ProjectItem{
		Name:            pinnedItem.Name,
		Articlelink:     articlelink,
		Codelink:        pinnedItem.URL,
		Demolink:        demoLink,
		Details:         pinnedItem.Description,
		DisplayPriority: sequenceString,
		Domain:          caser.String(strings.ToLower(domain)),
		Stack:           stack,
	}

	return projectItem
}

func getExistingValues(pinnedItem github.PinnedItem, existingProjects []model.ProjectItem) (string, string) {
	for _, p := range existingProjects {
		if strings.EqualFold(p.Name, pinnedItem.Name) {
			return p.Domain, p.Articlelink
		}
	}
	return "", ""
}

func enrichStack(pinnedItem github.PinnedItem) string {
	stackList := []string{}
	for _, item := range pinnedItem.Languages.Edges {
		stackList = append(stackList, item.Node.Name)
	}
	stack := strings.Join(stackList, ", ")

	return stack
}

func enrichDomain(pinnedItem github.PinnedItem) string {
	if len(pinnedItem.RepositoryTopics.Edges) == 0 {
		return domains["default"]
	}

	for _, item := range pinnedItem.RepositoryTopics.Edges {
		domain, ok := domains[strings.ToLower(item.Node.Topic.Name)]
		if ok {
			return domain
		}
	}

	return domains["default"]
}
