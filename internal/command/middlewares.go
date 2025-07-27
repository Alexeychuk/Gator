package command

import (
	"context"

	"github.com/Alexeychuk/Gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.Username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
