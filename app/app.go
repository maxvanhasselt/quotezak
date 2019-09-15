package app

import (
	"flag"
	"fmt"
	"os"

	"git.code-cloppers.com/max/quotezak/bot"
	"git.code-cloppers.com/max/quotezak/db"
	"git.code-cloppers.com/max/quotezak/messaging"
)

// Application is a runnable object that will run the service.
type Application struct {
	args []string
	db   *db.Database
}

// New returns an Application object that can be started with Run()
func New(args []string) *Application {
	return &Application{
		args: args,
	}
}

// Run starts the application
func (a *Application) Run() error {

	// Set the flag options
	var showVersion, showUsage, setupDb bool
	var envPath string
	fs := flag.NewFlagSet("quotezak", flag.ExitOnError)
	fs.SetOutput(os.Stdout)

	fs.BoolVar(&showUsage, "help", false, "Show this message")
	fs.BoolVar(&showVersion, "version", false, "Print version info")
	fs.BoolVar(&setupDb, "setup-db", false, "Set up the database for first use")
	fs.StringVar(&envPath, "env", ".env.yml", "path to env file")

	// Parse commandline arguments
	err := fs.Parse(a.args[1:])
	if err != nil {
		return err
	}

	if showUsage {
		fs.PrintDefaults()
		return nil
	}
	// Read the config file
	var cfg Config
	err = cfg.FromFile(envPath)
	if err != nil {
		return err
	}

	fmt.Print(cfg.ToString())
	db := &db.Database{}

	err = a.initDb(db, &cfg.Database)
	if err != nil {
		return err
	}

	if setupDb {
		err := db.SetupDatabase()
		if err != nil {
			return err
		}
	}
	messenger := messaging.NewMessenger(db.DB)
	bot := bot.New(&cfg.Bot, messenger)
	return bot.Start()
	// q := models.NewQuote("test", "test", "test", 2019)

	// err = q.Save(a.db.DB)
	// if err != nil {
	// 	return err
	// }

}

func (a *Application) initDb(db *db.Database, cfg *db.Config) error {
	err := db.InitDb(cfg)
	if err != nil {
		return err
	}
	a.db = db

	return nil
}
