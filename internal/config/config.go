package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func ReadConfig() *viper.Viper {
	v := viper.New()
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("unable to get working directory:", err)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println("unable to get user's home directory:", err)
	}
	v.AddConfigPath(cwd)
	v.AddConfigPath(filepath.Join(homedir, "/.config/grec"))
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}

type config struct {
	Token string `yaml:"token"`
}

func ParseConfig(v *viper.Viper) *config {
	c := &config{}
	err := v.Unmarshal(c)
	if err != nil {
		panic(err)
	}
	return c
}
