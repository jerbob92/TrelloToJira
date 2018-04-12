package helpers

import (
	"github.com/adlio/trello"
	"github.com/andygrunwald/go-jira"
)

type MigrateListMap struct {
	TrelloList trello.List
	JiraStatus string
}

type MigrateUserMap struct {
	TrelloMember trello.Member
	JiraUser     jira.User
}

type MigrateMap struct {
	TrelloBoard trello.Board
	JiraBoard   JiraProject
	ListMap     []MigrateListMap
	UserMap     []MigrateUserMap
}

type JiraProject struct {
	Expand          string               `json:"expand" structs:"expand"`
	Self            string               `json:"self" structs:"self"`
	ID              string               `json:"id" structs:"id"`
	Key             string               `json:"key" structs:"key"`
	Name            string               `json:"name" structs:"name"`
	AvatarUrls      jira.AvatarUrls      `json:"avatarUrls" structs:"avatarUrls"`
	ProjectTypeKey  string               `json:"projectTypeKey" structs:"projectTypeKey"`
	ProjectCategory jira.ProjectCategory `json:"projectCategory,omitempty" structs:"projectsCategory,omitempty"`
	IssueTypes      []jira.IssueType     `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
}
