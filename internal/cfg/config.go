package cfg

import (
	"github.com/d3code/xlog"
	"github.com/spf13/cobra"
)

type config struct {
	Verbose   bool
	Directory []string
	Extension []string
	Ignore    []string
	Command   []string
	Port      []int
}

var Config *config

func GetConfiguration(cmd *cobra.Command) {
	command := config{
		Verbose:   flagBool(cmd, "verbose"),
		Directory: flagStringSlice(cmd, "directory"),
		Extension: flagStringSlice(cmd, "extension"),
		Ignore:    flagStringSlice(cmd, "ignore"),
		Command:   flagStringSlice(cmd, "command"),
		Port:      flagIntSlice(cmd, "port"),
	}

	const prefix = "[RUN]"
	if command.Verbose {
		xlog.EnableConsole(xlog.LevelTrace, xlog.CallerShort, prefix, true)
	} else {
		xlog.EnableConsole(xlog.LevelInfo, xlog.CallerNone, prefix, true)
	}

	Config = &command
}

func flagBool(cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		xlog.Fatalf("Error parsing flag [ %s ]: %s", flag, err)
	}
	return value
}

func flagStringSlice(cmd *cobra.Command, flag string) []string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		xlog.Fatalf("Error parsing flag [ %s ]: %s", flag, err)
	}
	return value
}

func flagIntSlice(cmd *cobra.Command, flag string) []int {
	value, err := cmd.Flags().GetIntSlice(flag)
	if err != nil {
		xlog.Fatalf("Error parsing flag [ %s ]: %s", flag, err)
	}
	return value
}
