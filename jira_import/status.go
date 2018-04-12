package jira_import

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

type Status struct {
	ID   string `json:"id" structs:"id"`
	Name string `json:"name" structs:"name"`
	Self string `json:"self" structs:"self"`
}

func GetStatuses(client jira.Client) (*[]Status, *jira.Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/status")
	req, err := client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new([]Status)
	resp, err := client.Do(req, result)

	if err != nil {
		err = jira.NewJiraError(resp, err)
	}

	return result, resp, err
}
