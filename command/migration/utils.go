package migration

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"

	"github.com/gogovan-korea/ggx-kr-service-utils/command/constants"
	"github.com/gogovan-korea/ggx-kr-service-utils/config"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func MigrateDatabase(dbConfigs []databaseConfig, isUp bool, step int) {
	for _, cfg := range dbConfigs {
		execute(cfg, isUp, step)
	}
}
func execute(cfg databaseConfig, isUp bool, step int) {
	fmt.Println("Running migrate command")

	driver, err := getDbDriver(cfg)

	if err != nil {
		fmt.Println("Get db driver", zap.Error(err))
	}

	fileSource, err := (&file.File{}).Open(fmt.Sprintf("file://%s", cfg.MigrationFilePath))
	if err != nil {
		fmt.Println("opening file error", zap.Error(err))
	}

	m, err := migrate.NewWithInstance("file", fileSource, cfg.DbType, driver)
	if err != nil {
		fmt.Println("migrate error", zap.Error(err))
	}

	// Force if version exists
	version := viper.GetInt(constants.ForceFlagName)
	if version != 0 {
		err := m.Force(version)
		if err != nil {
			return
		}
	}

	if step == 0 {
		if isUp {
			err = m.Up()
		} else {
			err = m.Down()
		}
	} else {
		if isUp {
			err = m.Steps(step)
		} else {
			err = m.Steps(step * -1)
		}
	}

	if err == nil {
		fmt.Println("Migrate done with success")
	} else {
		if err.Error() != constants.NoChange {
			fmt.Println("migrate error", zap.Error(err))
			os.Exit(1)
		} else {
			fmt.Println("No change")
		}
	}
}

func getDbDriver(cfg databaseConfig) (database.Driver, error) {
	if cfg.DbType == "mysql" {
		db, err := sqlx.Connect(cfg.DbType, cfg.ConnectionString)
		if err != nil {
			panic(err)
		}
		return mysql.WithInstance(db.DB, &mysql.Config{})
	} else if cfg.DbType == "postgres" {
		db, err := sqlx.Connect(cfg.DbType, cfg.ConnectionString)
		if err != nil {
			panic(err)
		}
		return postgres.WithInstance(db.DB, &postgres.Config{})
	} else if cfg.DbType == "mongo" {
		client, err := mongo.NewClient(options.Client().ApplyURI(cfg.ConnectionString))
		if err != nil {
			panic(err)
		}
		if err := client.Connect(context.Background()); err != nil {
			return nil, err
		}

		return mongodb.WithInstance(client, &mongodb.Config{DatabaseName: cfg.Database})
	} else {
		return nil, errors.New("Unknown db type")
	}
}

func GetDbConfigs(dbConfigKeys ...string) []databaseConfig {
	var cfg interface{}
	config.LoadConfig(viper.GetString(constants.ConfigFlagName), &cfg)

	data := cfg.(map[string]interface{})

	var dbConfigs []databaseConfig
	for _, key := range dbConfigKeys {
		keys := strings.Split(key, ".")
		var cfgData map[string]interface{}
		for i, k := range keys {
			if i == 0 {
				cfgData = data[strings.ToLower(k)].(map[string]interface{})
			} else {
				cfgData = cfgData[strings.ToLower(k)].(map[string]interface{})
			}
		}
		dbName, _ := cfgData["database"].(string)

		dbConfig := databaseConfig{
			DbType:            cfgData["dbtype"].(string),
			ConnectionString:  cfgData["connectionstring"].(string),
			MigrationFilePath: cfgData["migrationfilepath"].(string),
			Database:          dbName,
		}
		dbConfigs = append(dbConfigs, dbConfig)
	}
	return dbConfigs
}
