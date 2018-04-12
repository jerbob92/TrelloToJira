package jira_import

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
)

// Find searches for user info from JIRA:
// It can find users by email, username or name
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-findUsers
func FindUsers(client jira.Client, project string) ([]jira.User, *jira.Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user/assignable/search?project=%s", project)
	req, err := client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	users := []jira.User{}
	resp, err := client.Do(req, &users)
	if err != nil {
		return nil, resp, jira.NewJiraError(resp, err)
	}
	return users, resp, nil
}
