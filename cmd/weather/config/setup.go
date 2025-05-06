package config

import (
	"os"

	"github.com/rs/zerolog/log"

	"sigs.k8s.io/yaml"
)

var Config tConfig

type tConfig struct {
	Weather tWeather
}

type tWeather struct {
	BaseUrl string `json:"baseUrl"`
	Key     []byte `json:"key"`
	Sub     string `json:"sub"`
	Alg     string `json:"alg"`
	Kid     string `json:"kid"`
}

func SetConfig(configDir string) {
	data, err := os.ReadFile(configDir + "/config/config.yaml")
	if err != nil {
		log.Panic().Err(err).Msg("")
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Panic().Err(err).Msg("SetConfig unmarshal")
	}

	// 需要自行创建和风天气的私钥文件
	keyData, err := os.ReadFile(configDir + "/config/private_key.pem")
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	Config.Weather.Key = keyData
}
