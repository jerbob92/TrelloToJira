# TrelloToJira [![Build Status](https://travis-ci.org/jerbob92/TrelloToJira.svg?branch=master)](https://travis-ci.org/jerbob92/TrelloToJira)
This is a quick'n'dirty Trello to Jira migration tool, created because the built in Jira tool did not work.

The tool is interactive, it will ask you questions before migrating.

## Capabilities
 - Map a Trello board to a Jira project
 - Map a Trello list (column) to a Jira status
 - Map a Trello member to a Jira user
 - Migrate cards
 - Migrate card attachments
 - Migrate card comments

## What's missing?
Probably some stuff that I didn't need, like checklists. The code is quite simple, MR's are welcome.

## Options
| Parameter              | Default Value    | Description  |
| ---------------------- | --------------   | ------------ |
| --trello_app_key              | ""        | Trello API App Key |
| --trello_token                 | ""               | Trello API token |
| --jira_url                  | ""               | Jira environment Base URL, with trailing slash |
| --jira_username          | ""               | Jira Username |
| --jira_password        | ""               | Jira Password |

## Builds
The package has several builds, for many different operating systems, check the releases tab.

If you do not wish to download a binary, and have go running locally, you can execute the following command:
```go install```