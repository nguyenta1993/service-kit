package command

import (
	"os"

	"github.com/tikivn/s14e-backend-utils/command/constants"
	"github.com/tikivn/s14e-backend-utils/command/migration"
	"github.com/tikivn/s14e-backend-utils/command/start"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func UseCommands(commands ...*cobra.Command) {
	var rootCmd = &cobra.Command{}

	pflag.String(constants.ConfigFlagName, "", "--config=<config-path>")
	pflag.Int(constants.ForceFlagName, 0, "--force=<version>")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func WithStartCommand(startFunc func(), cfg interface{}, dbConfigKeys ...string) *cobra.Command {
	return start.WithStartCommand(startFunc, cfg, dbConfigKeys...)
}

func WithMigrationCommand(dbConfigKeys ...string) *cobra.Command {
	return migration.MigrationCommand(dbConfigKeys...)
}
