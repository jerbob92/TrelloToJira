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

You can also set the options by using your environment. Use the same names but in uppercase.

## Running
Download any of the releases that matches your platform. You don't need any dependencies to run the program

## Compiling
I used Go 1.9 to develop, but it should work on most versions, make sure that you have [Godeps](https://github.com/tools/godep) available.
```
git clone https://github.com/jerbob92/TrelloToJira.git
cd TrelloToJira
godep restore
go build
./TrelloToJira
```

After the ```go build```, ```TrelloToJira``` will contain the compiled binary.
You can also use ```go run``` to run the program without compiling.
