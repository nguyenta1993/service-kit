package config

import (
	"github.com/gogovan-korea/ggx-kr-service-utils/vault"
	vaultgo "github.com/mittwald/vaultgo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func LoadConfig(configPath string, config interface{}) {
	if configPath == "" {
		panic("Missing config path")
	}

	viper.SetConfigType(Yaml)
	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	//Update config data if it's not the development env
	if viper.GetString("vault.address") != "" && !viper.GetBool("development") {
		updateDataFromVault()
	}

	if config != nil {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	}
}

func updateDataFromVault() {
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
