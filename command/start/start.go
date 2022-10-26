package start

import (
	"github.com/tikivn/s14e-backend-utils/command/constants"
	"github.com/tikivn/s14e-backend-utils/command/migration"
	"github.com/tikivn/s14e-backend-utils/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithStartCommand(startFunc func(), cfg interface{}, dbConfigKeys ...string) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the server",
		Run: func(cmd *cobra.Command, args []string) {
			config.LoadConfig(viper.GetString(constants.ConfigFlagName), cfg)

			if len(dbConfigKeys) != 0 {
				dbConfigs := migration.GetDbConfigs(dbConfigKeys...)
				migration.MigrateDatabase(dbConfigs, true, 0)
			}

			startFunc()
		},
	}
}
