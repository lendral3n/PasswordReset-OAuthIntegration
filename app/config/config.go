package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	JWT_SECRET string
	RDS_URL string
	AWS_ACCESS_KEY_ID string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION string
)

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOSTNAME string
	DB_PORT     int
	DB_NAME     string
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("DBUSER"); found {
		app.DB_USERNAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASS"); found {
		app.DB_PASSWORD = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DB_HOSTNAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, _ := strconv.Atoi(val)
		app.DB_PORT = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DB_NAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("JWTSECRET"); found {
		JWT_SECRET = val
		isRead = false
	}
	if val, found := os.LookupEnv("RDSURL"); found {
		RDS_URL = val
		isRead = false
	}
	if val, found := os.LookupEnv("AWSKEYID"); found {
		AWS_ACCESS_KEY_ID = val
		isRead = false
	}
	if val, found := os.LookupEnv("AWSSECRET"); found {
		AWS_SECRET_ACCESS_KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("AWSREGION"); found {
		AWS_REGION = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}

		AWS_ACCESS_KEY_ID = viper.GetString("AWSKEY")
		AWS_SECRET_ACCESS_KEY = viper.GetString("AWSSECRET")
		AWS_REGION = viper.GetString("AWSREGION")
		RDS_URL = viper.GetString("RDSURL")
		JWT_SECRET = viper.GetString("JWTSECRET")
		app.DB_USERNAME = viper.Get("DBUSER").(string)
		app.DB_PASSWORD = viper.Get("DBPASS").(string)
		app.DB_HOSTNAME = viper.Get("DBHOST").(string)
		app.DB_PORT, _ = strconv.Atoi(viper.Get("DBPORT").(string))
		app.DB_NAME = viper.Get("DBNAME").(string)
	}

	return &app
}
