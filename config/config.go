package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabaseName                string   `json:"database_name"`
	DatabaseConfig              string   `json:"database_config"`
	Debug                       bool     `json:"debug"`
	MicrosoftOAUTH2ClientID     string   `json:"ms_oauth2_client_id"`
	MicrosoftOAUTH2Secret       string   `json:"ms_oauth2_secret"`
	MicrosoftOAUTH2RefreshToken string   `json:"ms_oauth2_refresh_token"`
	Webhooks                    []string `json:"webhooks"`
}

func GetConfig() (Config, error) {
	var config Config
	file, err := os.ReadFile("config.json")
	if err != nil {
		marshal, err := json.Marshal(Config{
			DatabaseName:                "sqlite3",
			DatabaseConfig:              "database.sqlite3",
			Debug:                       true,
			MicrosoftOAUTH2ClientID:     "",
			MicrosoftOAUTH2RefreshToken: "",
			MicrosoftOAUTH2Secret:       "",
			Webhooks:                    make([]string, 0),
		})
		if err != nil {
			return config, err
		}
		err = os.WriteFile("config.json", marshal, 0600)
		if err != nil {
			return config, err
		}
		file, err = os.ReadFile("config.json")
		if err != nil {
			return config, err
		}
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

func SaveConfig(config Config) error {
	marshal, err := json.Marshal(config)
	if err != nil {
		return err
	}
	f, err := os.Create("config.json")
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	err = os.WriteFile("config.json", marshal, 0600)
	if err != nil {
		return err
	}
	return nil
}
