package app

import (
	"flag"
	"fmt"
	"os"
)

// Application is a runnable object that will run the service.
type Application struct {
	args []string
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
	var showVersion, showUsage bool
	var envPath string
	fs := flag.NewFlagSet("quotezak", flag.ExitOnError)
	fs.SetOutput(os.Stdout)

	fs.BoolVar(&showUsage, "help", false, "Show this message")
	fs.BoolVar(&showVersion, "version", false, "Print version info")
	fs.StringVar(&envPath, "env", ".env.yml", "path to env file")

	// Parse commandline arguments
	err := fs.Parse(a.args[1:])
	if err != nil {
		return err
	}

	fmt.Print(envPath)
	// Read the config file
	var cfg Config
	err = cfg.FromFile(envPath)
	if err != nil {
		return err
	}

	fmt.Print(cfg.ToString())

	return nil
}
