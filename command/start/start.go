package start

import (
	"fmt"
	"github.com/gogovan-korea/ggx-kr-service-utils/command/constants"
	"github.com/gogovan-korea/ggx-kr-service-utils/command/migration"
	"github.com/gogovan-korea/ggx-kr-service-utils/config"
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
