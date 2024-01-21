package cli

import (
	"flag"

	"github.com/brecht-vde/reading-list/internal"
)

type createDatabaseCommand struct {
	ParentId string
	Secret   string
}

func NewCreateDateBaseCommand(args []string) (createDatabaseCommand, error) {
	var cmd createDatabaseCommand
	var err error

	var id = flag.String("id", "", "the notion page id to which a database will be added")
	var secret = flag.String("secret", "", "the notion internal integration secret")

	err = flag.CommandLine.Parse(args)

	if err != nil {
		return cmd, err
	}

	cmd = createDatabaseCommand{
		ParentId: *id,
		Secret:   *secret,
	}

	return cmd, err
}

func (c *createDatabaseCommand) Run() error {
	var err error

	client := internal.NewNotionClient(c.Secret)
	parameters := internal.CreateDatabaseRequestParameters{
		ParentId: c.ParentId,
		Title:    "Reading list",
		Url:      "github.com/brecht-vde/reading-list",
	}
	body := internal.NewCreateDatabaseRequest(parameters)

	err = client.Post("databases", body)

	return err
}
