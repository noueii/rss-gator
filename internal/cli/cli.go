package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/noueii/rss-gator/internal/api"
	"github.com/noueii/rss-gator/internal/app"
	"github.com/noueii/rss-gator/internal/db"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*app.App, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*app.App, Command) error),
	}

}

func (c *Commands) Register(name string, f func(*app.App, Command) error) {
	c.handlers[name] = f
}

func (c *Commands) Run(a *app.App, cmd Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}

	return handler(a, cmd)
}

func HandlerLogin(a *app.App, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required")
	}

	username := cmd.Args[0]

	_, err := a.DB.GetUserByName(context.Background(), username)

	if err != nil {
		return err
	}

	if err := a.Config.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", username)

	return nil
}

func HandlerRegister(a *app.App, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required")
	}

	username := cmd.Args[0]
	createdUser, err := a.DB.CreateUser(context.Background(), username)

	if err != nil {
		fmt.Println("Could not create user")
		return err
	}

	fmt.Println("User created successfully")
	printUser(createdUser)

	return a.Config.SetUser(createdUser.Name)
}

func printUser(user db.User) {
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Printf("Updated at: %s\n", user.UpdatedAt)
}

func HandlerReset(a *app.App, cmd Command) error {
	if err := a.DB.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func HandlerUsers(a *app.App, cmd Command) error {
	users, err := a.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	printUsers(users, a.Config.CurrentUsername)
	return nil
}

func printUsers(users []db.User, currentUsername string) {
	for _, user := range users {
		isCurrent := user.Name == currentUsername
		str := fmt.Sprintf("* %s ", user.Name)
		if isCurrent {
			str += "(current)\n"
		} else {
			str += "\n"
		}
		fmt.Print(str)
	}

}

func HandleAgg(a *app.App, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("ERROR: The command should contain the interval between the updates.\nUse 's' for seconds, 'm' for minutes, 'h' for hours etc.\nExample: agg 10s ; agg 10m ; agg 10h\n")
	}

	duration, err := time.ParseDuration(cmd.Args[0])

	fmt.Printf("Collecting feeds every %s\n", duration.String())

	if err != nil {
		return err
	}

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		if err := api.ScrapeFeeds(a); err != nil {
			return err
		}
	}

	return nil
}

func HandleAddFeed(a *app.App, cmd Command, user db.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Too few arguments. Please use it as: 'addfeed [name] [link]'")
	}

	createdFeed, err := a.DB.CreateFeed(context.Background(), db.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	})

	fmt.Printf("createdFeed: %v\n", createdFeed)

	if err != nil {
		return err
	}
	return HandleFollow(a, Command{
		Name: "follow",
		Args: []string{createdFeed.Url},
	}, user)

}

func HandleFeeds(a *app.App, cmd Command) error {
	feeds, err := a.DB.GetFeedsWithAuthor(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s\tURL:%s\tAuthor:%s\n", feed.Name, feed.Url, feed.Name_2)
	}

	return nil
}

func HandleFollow(a *app.App, cmd Command, user db.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Your follow command should only contain the url of the feed. Example: follow [url]")
	}
	feedURL := cmd.Args[0]
	feed, err := a.DB.GetFeedByURL(context.Background(), feedURL)

	if err != nil {
		return err
	}

	feedFollow, err := a.DB.CreateFeedFollow(context.Background(), db.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("User %s followed feed %s", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func HandleFollowing(a *app.App, cmd Command, user db.User) error {
	feeds, err := a.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Println("You are not following any feed")
	} else {
		fmt.Println("You are currently following: ")
	}

	for _, feed := range feeds {
		fmt.Printf("\t- %s\n", feed.FeedName)
	}

	return nil
}

func HandleUnfollow(a *app.App, cmd Command, user db.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Your command should only contain the url of the feed. Example unfollow [url]")
	}

	feedUrl := cmd.Args[0]

	feed, err := a.DB.GetFeedByURL(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	if err := a.DB.DeleteUserFeed(context.Background(), db.DeleteUserFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return err
	}
	return nil

}

func HandleBrowse(a *app.App, cmd Command, user db.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		integer, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = integer
	}

	posts, err := a.DB.GetPostsForUserWithLimit(context.Background(), db.GetPostsForUserWithLimitParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return err
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post db.Post) {
	fmt.Printf("TITLE: %s\n", post.Title)
	fmt.Printf("PUBLICATION DATE: %s\n", post.PublishedAt)
	fmt.Printf("LINK: %s\n", post.Url)
	fmt.Printf("DESCRIPTION: %s\n\n", post.Description)
}
