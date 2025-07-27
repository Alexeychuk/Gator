package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Alexeychuk/Gator/internal/config"
	"github.com/Alexeychuk/Gator/internal/database"
	rssfeed "github.com/Alexeychuk/Gator/internal/rssFeed"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commandsMap map[string]func(s *State, cmd Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {

	if s.Cfg == nil {
		return errors.New("no config")
	}

	command, exists := c.commandsMap[cmd.Name]

	if !exists {
		return errors.New("no command found")
	}

	err := command(s, cmd)

	if err == nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.commandsMap == nil {
		c.commandsMap = make(map[string]func(s *State, cmd Command) error)
	}

	c.commandsMap[name] = f
}

// handlers

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Print("the login handler expects a single argument, the username\n")
		os.Exit(1)
	}

	user, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Print("user doesnt exists in db\n")
		os.Exit(1)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set\n")

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Print("the register handler expects a single argument, the username\n")
		os.Exit(1)
	}
	name := cmd.Args[0]

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})

	if err != nil {
		// Check if it's a PostgreSQL error
		if pqErr, ok := err.(*pq.Error); ok {
			// Check for unique violation error code
			if pqErr.Code == "23505" {
				fmt.Printf("user %s already exists", name)
				os.Exit(1)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Print("user table reset error\n")
		fmt.Println(err)

		os.Exit(1)
	}

	fmt.Printf("User table has been reset\n")

	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Print("error in GetUsers\n")
		os.Exit(1)
	}

	for _, user := range users {

		if s.Cfg.Username == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	// if len(cmd.Args) == 0 {
	// 	fmt.Print("the register handler expects a single argument, the username\n")
	// 	os.Exit(1)
	// }
	// feedUrl := cmd.Args[0]

	data, err := rssfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Print("error in GetUsers\n")
		os.Exit(1)
	}

	fmt.Println(data)

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		fmt.Print("the register handler expects a two arguments, the name and urlf\n")
		os.Exit(1)
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{ID: uuid.New(), Name: cmd.Args[0], Url: cmd.Args[1], UserID: user.ID})
	if err != nil {
		fmt.Print("failed to create feed due to next error:\n")
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})

	if err != nil {
		fmt.Print("failed to create feed_follow due to next error:\n")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(feed)
	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("Cant find user for feed: %s", feed.Name)
			return err
		}

		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", user.Name)

	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Print("the follow handler expects one arguments, url\n")
		os.Exit(1)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	feed_follow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}

	fmt.Printf("Follow - feed name: %s\nuser_name: %s\n", feed_follow.FeedName, feed_follow.UserName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {

	user_feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Feeds followed by user %s\n", user.Name)
	for _, user_feed := range user_feeds {
		fmt.Printf("- %s\n", user_feed.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		fmt.Print("the unfollow handler expects one arguments, url\n")
		os.Exit(1)
	}

	err := s.Db.DeleteFeedFollowByUserIdAndUrl(context.Background(), database.DeleteFeedFollowByUserIdAndUrlParams{UserID: user.ID, Url: cmd.Args[0]})
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s unfollowed by user %s\n", cmd.Args[0], user.Name)

	return nil
}
