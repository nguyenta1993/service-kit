package config

import (
	vaultgo "github.com/mittwald/vaultgo"
	"github.com/spf13/viper"
	"github.com/gogovan-korea/ggx-kr-service-utils/vault"
)

func LoadConfig(configPath string, config interface{}) {
	if configPath == "" {
		panic("Missing config path")
	}

	viper.SetConfigType(Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Update config data if it's not the development env
	if viper.GetString("vault.address") != "" {
		updateDataFromVault()
	}

	viper.AutomaticEnv()

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

	var vaultClient *vault.VaultClient
	if token != "" {
		vaultClient, _ = vault.NewVaultClient(address, vaultgo.WithAuthToken(token))
	} else {
		vaultClient, _ = vault.NewVaultClient(address, vaultgo.WithKubernetesAuth(role, vaultgo.WithMountPoint(mountPoint)))
	}

	secretData, _ := vaultClient.GetSecretKeys(path)

	for key, s := range secretData {
		viper.Set(key, s)
	}
}
