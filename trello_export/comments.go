package trello_export

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/adlio/trello"
)

func GetComments(client *trello.Client, c *trello.Card, args trello.Arguments) ([]trello.Action, error) {
	path := fmt.Sprintf("cards/%s/actions", c.ID)
	args["filter"] = "commentCard"
	action := []trello.Action{}
	err := client.Get(path, args, &action)
	if err != nil {
		err = errors.Wrapf(err, "Error getting comments on card %s", c.ID)
	}
	return action, err
}