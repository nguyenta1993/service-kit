package config

import (
	"github.com/gogovan/ggx-kr-service-utils/vault"
	vaultgo "github.com/mittwald/vaultgo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func LoadConfig(configPath string, config interface{}) {
	if configPath == "" {
		panic("Missing config path")
	}

	localViper := viper.New()
	localViper.SetConfigType(Yaml)
	localViper.SetConfigFile(configPath)
	localViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	localViper.AutomaticEnv()
	if err := localViper.ReadInConfig(); err != nil {
		panic(err)
	}

	// update data from .vaultenv file
	if !localViper.GetBool("development") {
		vaultEnvPath := os.Getenv("VAULT_ENV_PATH")
		if vaultEnvPath == "" {
			vaultEnvPath = "/app/.vaultenv"
		}
		println("Load config from external vault: " + vaultEnvPath + "")
		localViper.SetConfigFile(vaultEnvPath)
		localViper.SetConfigType(Env)
		if err := localViper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				if viper.GetString("vault.address") != "" {
					updateDataFromVault()
				}
			}

		}
	}
	if config != nil {
		if err := localViper.Unmarshal(config); err != nil {
			panic(err)
		}
	}
}

func updateDataFromVault() {
	println("Fall back to update config from remote vault")
	address := viper.GetString("vault.address")
	path := viper.GetString("vault.path")
	token := viper.GetString("vault.token")
	role := viper.GetString("vault.role")
	mountPoint := viper.GetString("vault.mountPoint")

	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	var vaultClient *vault.VaultClient
	if token != "" {
		vaultClient, _ = vault.NewVaultClient(address, vaultgo.WithAuthToken(token))
	} else {
		vaultClient, _ = vault.NewVaultClient(address, vaultgo.WithKubernetesAuth(role, vaultgo.WithMountPoint(mountPoint)))
	}

	secretData, err := vaultClient.GetSecretKeys(path)
	if err != nil {
		panic(err)
	}
	for key, s := range secretData {
		viper.Set(key, s)
	}
}
