package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Alexeychuk/Gator/internal/command"
	"github.com/Alexeychuk/Gator/internal/config"
	"github.com/Alexeychuk/Gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	foundConfig, err := config.Read()
	if err != nil {
		return
	}

	state := command.State{
		Cfg: foundConfig,
	}

	commands := command.Commands{}

	commands.Register("login", command.HandlerLogin)
	commands.Register("register", command.HandlerRegister)
	commands.Register("reset", command.HandlerReset)
	commands.Register("users", command.HandlerGetUsers)
	commands.Register("agg", command.HandlerAgg)
	commands.Register("addfeed", command.MiddlewareLoggedIn(command.HandlerAddFeed))
	commands.Register("feeds", command.HandlerGetFeeds)
	commands.Register("follow", command.MiddlewareLoggedIn(command.HandlerFollow))
	commands.Register("following", command.MiddlewareLoggedIn(command.HandlerFollowing))
	commands.Register("unfollow", command.MiddlewareLoggedIn(command.HandlerUnfollow))

	db, err := sql.Open("postgres", foundConfig.DBUrl)

	if err != nil {
		fmt.Println(err)
		return
	}

	dbQueries := database.New(db)
	state.Db = dbQueries

	if len(os.Args) < 2 {
		fmt.Print("Not enough arguments\n")
		os.Exit(1)
	}

	// command running

	var madeCommand command.Command
	madeCommand.Name = os.Args[1]
	madeCommand.Args = os.Args[2:]
	fmt.Println(os.Args)

	err = commands.Run(&state, madeCommand)

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("db: %s, user: %s\n", foundConfig.DBUrl, foundConfig.Username)

}
