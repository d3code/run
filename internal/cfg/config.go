package cfg

import (
	"github.com/d3code/xlog"
	"github.com/spf13/cobra"
	"os"
)

type config struct {
	Verbose   bool
	Directory []string
	Extension []string
	Ignore    []string
	Run       []string
	Port      []int
}

var Config *config

func GetConfiguration(cmd *cobra.Command) {
	command := config{
		Verbose:   flagBool(cmd, "verbose"),
		Directory: flagStringSlice(cmd, "directory"),
		Extension: flagStringSlice(cmd, "extension"),
		Ignore:    flagStringSlice(cmd, "ignore"),
		Run:       flagStringSlice(cmd, "run"),
		Port:      flagIntSlice(cmd, "port"),
	}

	if command.Verbose {
		xlog.EnableConsole(xlog.LevelTrace, xlog.CallerShort, true)
		xlog.Info("Verbose mode enabled")
	}

	Config = &command
}

func flagBool(cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		xlog.Errorf("error parsing flag [ %s ]: %s", flag, err.Error())
		os.Exit(1)
	}
	return value
}

func flagStringSlice(cmd *cobra.Command, flag string) []string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		xlog.Errorf("error parsing flag [ %s ]: %s", flag, err.Error())
		os.Exit(1)
	}
	return value
}

func flagIntSlice(cmd *cobra.Command, flag string) []int {
	value, err := cmd.Flags().GetIntSlice(flag)
	if err != nil {
		xlog.Errorf("error parsing flag [ %s ]: %s", flag, err.Error())
		os.Exit(1)
	}
	return value
}
