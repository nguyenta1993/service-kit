package start

import (
	"fmt"
	"github.com/nguyenta1993/service-kit/command/constants"
	"github.com/nguyenta1993/service-kit/command/migration"
	"github.com/nguyenta1993/service-kit/config"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithStartCommand(startFunc func(), cfg interface{}, dbConfigKeys ...string) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the server",
		Run: func(cmd *cobra.Command, args []string) {
			var configPath string
			// Priority config from env
			if environ := os.Getenv(config.AppEnv); environ != "" {
				configPath = fmt.Sprintf("./config/%s/config.yaml", environ)
			} else {
				configPath = viper.GetString(constants.ConfigFlagName)
			}
			config.LoadConfig(configPath, cfg)
			if len(dbConfigKeys) != 0 {
				dbConfigs := migration.GetDbConfigs(dbConfigKeys...)
				migration.MigrateDatabase(dbConfigs, true, 0)
			}

			startFunc()
		},
	}
}
