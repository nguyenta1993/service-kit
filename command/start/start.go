package start

import (
	"github.com/gogovan-korea/ggx-kr-service-utils/command/constants"
	"github.com/gogovan-korea/ggx-kr-service-utils/command/migration"
	"github.com/gogovan-korea/ggx-kr-service-utils/config"

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
