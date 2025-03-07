package main

import (
	"fmt"
	"os"

	"github.com/noueii/rss-gator/internal/app"
	"github.com/noueii/rss-gator/internal/cli"
	"github.com/noueii/rss-gator/internal/config"
	"github.com/noueii/rss-gator/internal/middleware"
)

type State struct {
	config *config.Config
}

func main() {
	app, err := app.New()

	if err != nil {
		fmt.Println("Error initializing application: ", err)
		os.Exit(1)
	}

	defer app.Close()

	commands := cli.NewCommands()
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerUsers)
	commands.Register("agg", cli.HandleAgg)
	commands.Register("addfeed", middleware.LoggedIn(cli.HandleAddFeed))
	commands.Register("feeds", cli.HandleFeeds)
	commands.Register("follow", middleware.LoggedIn(cli.HandleFollow))
	commands.Register("following", middleware.LoggedIn(cli.HandleFollowing))
	commands.Register("unfollow", middleware.LoggedIn(cli.HandleUnfollow))
	commands.Register("browse", middleware.LoggedIn(cli.HandleBrowse))
	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(app, cmd); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
