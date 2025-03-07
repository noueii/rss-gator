package middleware

import (
	"context"

	"github.com/noueii/rss-gator/internal/app"
	"github.com/noueii/rss-gator/internal/cli"
	"github.com/noueii/rss-gator/internal/db"
)

func LoggedIn(handler func(a *app.App, c cli.Command, u db.User) error) func(*app.App, cli.Command) error {
	return func(a *app.App, cmd cli.Command) error {
		user, err := a.DB.GetUserByName(context.Background(), a.Config.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(a, cmd, user)
	}
}
