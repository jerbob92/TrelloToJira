package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jerbob92/TrelloToJira/helpers"
	"github.com/jerbob92/TrelloToJira/jira_import"

	"github.com/adlio/trello"
	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/jerbob92/TrelloToJira/trello_export"
	"time"
)

var TrelloClient *trello.Client
var JiraClient *jira.Client

func main() {
	TrelloAppKey := flag.String("trello_app_key", os.Getenv("TRELLO_APP_KEY"), "The Trello App Key.")
	TrelloToken := flag.String("trello_token", os.Getenv("TRELLO_TOKEN"), "The Trello Token.")
	JiraURL := flag.String("jira_url", os.Getenv("JIRA_URL"), "The JIRA endpoint")
	JiraUsername := flag.String("jira_username", os.Getenv("JIRA_USERNAME"), "The JIRA password")
	JiraPassword := flag.String("jira_password", os.Getenv("JIRA_PASSWORD"), "The JIRA username")

	if TrelloAppKey == nil || *TrelloAppKey == "" {
		log.Fatal("Missing Trello App Key")
	}

	if TrelloToken == nil || *TrelloToken == "" {
		log.Fatal("Missing Trello token")
	}

	if JiraURL == nil || *JiraURL == "" {
		log.Fatal("Missing JIRA URL")
	}

	if JiraUsername == nil || *JiraUsername == "" {
		log.Fatal("Missing JIRA Username")
	}

	if JiraPassword == nil || *JiraPassword == "" {
		log.Fatal("Missing JIRA Password")
	}

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	client := trello.NewClient(*TrelloAppKey, *TrelloToken)
	//client.Logger = logger

	TrelloClient = client

	tp := jira.CookieAuthTransport{
		Username: *JiraUsername,
		Password: *JiraPassword,
		AuthURL:  *JiraURL + "rest/auth/1/session",
	}

	jira_client, err := jira.NewClient(tp.Client(), *JiraURL)
	if err != nil {
		log.Fatalf("Could not get JIRA client: %s", err.Error())
	}

	JiraClient = jira_client

	JiraClient.Authentication.AcquireSessionCookie(*JiraUsername, *JiraPassword)

	current_user, err := JiraClient.Authentication.GetCurrentUser()
	if err != nil {
		log.Fatalf("Could not get JIRA client: %s", err.Error())
	}

	log.Printf("Logged in at Jira for user %s", current_user.Name)

	configureMigrate()
}

func configureMigrate() {
	token, err := TrelloClient.GetToken(TrelloClient.Token, trello.Defaults())
	if err != nil {
		log.Fatalf("Could not get details of Trello Token %s", err.Error())
	}

	trello_member, err := TrelloClient.GetMember(token.IDMember, trello.Defaults())
	if err != nil {
		log.Fatalf("Could not get member of Trello Token %s", err.Error())
	}

	trello_boards, err := trello_member.GetBoards(trello.Defaults())
	if err != nil {
		log.Fatalf("Could not get boards of Trello Token %s", err.Error())
	}

	if len(trello_boards) == 0 {
		log.Fatal("No Trello boards found")
	}

	jira_projects, _, err := JiraClient.Project.GetList()
	if err != nil {
		log.Fatalf("Could not get projects from Jira %s", err.Error())
	}

	if len(*jira_projects) == 0 {
		log.Fatal("No Jira projects found")
	}

	jira_project_options_display := []string{}
	jira_project_options := []string{}
	jira_project_list := map[string]helpers.JiraProject{}
	for _, jira_project := range *jira_projects {
		jira_project_list[jira_project.ID] = helpers.JiraProject{
			Expand:          jira_project.Expand,
			Self:            jira_project.Self,
			ID:              jira_project.ID,
			Key:             jira_project.Key,
			Name:            jira_project.Name,
			AvatarUrls:      jira_project.AvatarUrls,
			ProjectTypeKey:  jira_project.ProjectTypeKey,
			ProjectCategory: jira_project.ProjectCategory,
			IssueTypes:      jira_project.IssueTypes,
		}
		jira_project_options = append(jira_project_options, jira_project.ID)
		jira_project_options_display = append(jira_project_options_display, jira_project.ID+": "+jira_project.Key+" - "+jira_project.Name)
	}

	boardsToImport := []helpers.MigrateMap{}
	for _, trello_board := range trello_boards {
		default_answer := false
		if helpers.AskForConfirmation(fmt.Sprintf("Would you like to import board \"%s\"? [(Y)es/(N)o], default: No", trello_board.Name), &default_answer) {
			fmt.Printf("Marked board \"%s\" for import\n", trello_board.Name)
			selected_jira_project := helpers.AskForSelect(fmt.Sprintf("To which Jira project would you like to export this Trello board? [%s]", strings.Join(jira_project_options_display, " / ")), jira_project_options)

			trello_lists, err := trello_board.GetLists(trello.Defaults())
			if err != nil {
				log.Fatalf("Could not get lists from Trello %s", err.Error())
			}

			jira_project := jira_project_list[selected_jira_project]

			jira_statuses, _, err := jira_import.GetStatuses(*JiraClient)
			if err != nil {
				log.Fatalf("Could not get statuses from Jira %s", err.Error())
			}

			jira_column_options_display := []string{}
			jira_column_options := []string{}
			jira_column_list := map[string]string{}

			// Add option to skip column.
			jira_column_list["0"] = "0"
			jira_column_options = append(jira_column_options, strconv.Itoa(0))
			jira_column_options_display = append(jira_column_options_display, strconv.Itoa(0)+": Skip")

			for _, jira_status := range *jira_statuses {
				jira_column_list[jira_status.ID] = jira_status.ID
				jira_column_options = append(jira_column_options, jira_status.ID)
				jira_column_options_display = append(jira_column_options_display, jira_status.ID+": "+jira_status.Name)
			}

			importConfig := helpers.MigrateMap{
				TrelloBoard: *trello_board,
				JiraBoard:   jira_project,
			}

			for _, trello_list := range trello_lists {
				selected_jira_status := helpers.AskForSelect(fmt.Sprintf("To which Jira status would you like to export the Trello column \"%s\"? [%s]", trello_list.Name, strings.Join(jira_column_options_display, " / ")), jira_column_options)

				if selected_jira_status != "0" {
					importConfig.ListMap = append(importConfig.ListMap, helpers.MigrateListMap{
						TrelloList: *trello_list,
						JiraStatus: jira_column_list[selected_jira_status],
					})
				}
			}

			jira_users, _, err := jira_import.FindUsers(*JiraClient, jira_project.Key)
			if err != nil {
				log.Fatalf("Could not get users from Jira %s", err.Error())
			}

			if len(jira_users) == 0 {
				log.Fatal("No Jira users found")
			}

			jira_user_options_display := []string{}
			jira_user_options := []string{}
			jira_user_list := map[string]jira.User{}

			// Add option to skip column.
			jira_user_list["0"] = jira.User{
				Name: "0",
			}
			jira_user_options = append(jira_user_options, "0")
			jira_user_options_display = append(jira_user_options_display, strconv.Itoa(0)+": Unassigned")

			for _, jira_user := range jira_users {

				// Skip internal users.
				if strings.HasPrefix(jira_user.Name, "addon_") {
					continue
				}

				jira_user_list[jira_user.Name] = jira_user
				jira_user_options = append(jira_user_options, jira_user.Name)
				jira_user_options_display = append(jira_user_options_display, jira_user.Name+": "+jira_user.EmailAddress+" - "+jira_user.DisplayName)
			}

			trello_members, err := trello_board.GetMembers(trello.Defaults())
			if err != nil {
				log.Fatalf("Could not get members from Trello %s", err.Error())
			}

			for _, trello_member := range trello_members {
				selected_jira_member := helpers.AskForSelect(fmt.Sprintf("To which Jira member would you like to map the Trello member \"%s\"? [%s]", trello_member.Username+" - "+trello_member.FullName+" - "+trello_member.FullName, strings.Join(jira_user_options_display, " / ")), jira_user_options)

				if selected_jira_member != "0" {
					importConfig.UserMap = append(importConfig.UserMap, helpers.MigrateUserMap{
						TrelloMember: *trello_member,
						JiraUser:     jira_user_list[selected_jira_member],
					})
				}
			}

			boardsToImport = append(boardsToImport, importConfig)
		} else {
			fmt.Printf("Skipping board \"%s\"\n", trello_board.Name)
		}
	}

	executeMigrate(boardsToImport)
}

func executeMigrate(boards []helpers.MigrateMap) {
	if len(boards) == 0 {
		log.Fatalf("No boards selected for migrate")
	}

	for _, board := range boards {
		err := exportTrelloBoard(board)
		if err != nil {
			log.Fatalf("Could not migrate from Trello to Jira: %s", err.Error())
		}
	}
}

func exportTrelloBoard(board helpers.MigrateMap) error {
	for _, list_map := range board.ListMap {
		cards, err := list_map.TrelloList.GetCards(trello.Arguments{
			"attachments": "true",
		})
		if err != nil {
			return nil
		}

		log.Infof("Loaded %d cards for Trello board %s, list %s", len(cards), board.TrelloBoard.Name, list_map.TrelloList.Name)

		for _, card := range cards {
			status := &jira.Status{
				ID: list_map.JiraStatus,
			}

			var assignee *jira.User
			if card.Members != nil && len(card.Members) > 0 {
				for _, user_map := range board.UserMap {
					if user_map.TrelloMember.ID == card.Members[0].ID {
						assignee = &jira.User{
							Name: user_map.JiraUser.Name,
						}
						break
					}
				}
			}

			i := jira.Issue{
				Fields: &jira.IssueFields{
					Assignee:    assignee,
					Description: "Imported from Trello " + card.ShortUrl + "\n" + card.Desc,
					Type: jira.IssueType{
						Name: "Task",
					},
					Project: jira.Project{
						ID: board.JiraBoard.ID,
					},
					Summary: card.Name,
				},
			}

			issue, _, err := JiraClient.Issue.Create(&i)
			if err != nil {
				return err
			}

			log.Infof("Created issue %s", issue.Key)

			// Do transition.
			if status != nil {
				transitions, _, err := JiraClient.Issue.GetTransitions(issue.Key)
				if err != nil {
					return err
				}

				for _, transition := range transitions {
					if transition.To.ID == status.ID {
						_, err = JiraClient.Issue.DoTransition(issue.Key, transition.ID)
						if err != nil {
							log.Warnf("Warning: could not set transition for new issue %s: %s", issue.Key, err.Error())
						} else {
							log.Infof("Transitioned issue %s to status %s", issue.Key, status.Name)
						}
						break
					}
				}
			}

			if len(card.Attachments) > 0 {
				for _, attachment := range card.Attachments {
					resp, err := http.Get(attachment.URL)
					if err != nil {
						log.Warnf("Warning: could not download attachment %s: %s", attachment.Name, err.Error())

						// Skip file
						continue
					} else {
						log.Infof("Downloaded file %s for issue %s", attachment.Name, issue.Key)
					}

					_, _, err = JiraClient.Issue.PostAttachment(issue.Key, resp.Body, attachment.Name)
					if err != nil {
						log.Warnf("Warning: could not upload attachment %s: %s", attachment.Name, err.Error())
					} else {
						log.Infof("Uploaded file %s for issue %s", attachment.Name, issue.Key)
					}

					resp.Body.Close()
				}
			}

			card_comments, err := trello_export.GetComments(TrelloClient, card, trello.Defaults())
			if err != nil {
				log.Warnf("Warning: could not get comments for card %s: %s", card.Name, err.Error())
			} else {
				if len(card_comments) > 0 {
					for _, card_comment := range card_comments {
						jira_comment := jira.Comment{
							Body: "Imported from Trello " + card.ShortUrl + "\n" + card_comment.Data.Text,
							Created: card_comment.Date.Format(time.RFC3339),
						}

						if card_comment.MemberCreator != nil {
							for _, user_map := range board.UserMap {
								if user_map.TrelloMember.ID == card_comment.MemberCreator.ID {
									jira_comment.Author = jira.User{
										Name: user_map.JiraUser.Name,
									}
									break
								}
							}
						}

						_, _, err = JiraClient.Issue.AddComment(issue.Key, &jira_comment)
						if err != nil {
							log.Warnf("Warning: could not add comment for issue %s: %s", issue.Key, err.Error())
						} else {
							log.Infof("Added comment for issue %s", issue.Key)
						}
					}
				}
			}
		}
	}

	return nil
}
