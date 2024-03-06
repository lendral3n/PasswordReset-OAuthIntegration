package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var (
	JWT_SECRET            string
	RDS_URL               string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	CLIENT_ID             string
	CLIENT_SECRET         string
	GOOGLE_URL            string
	SCOPES                []string
	CLIENT_ID_FB          string
	CLIENT_SECRET_FB      string
	FB_URL                string
	SCOPES_FB             []string
)

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOSTNAME string
	DB_PORT     int
	DB_NAME     string
	SMTP_HOST   string
	SMTP_PORT   int
	SMTP_USER   string
	SMTP_PASS   string
	PASSWD_URL  string
	EMAIL_FROM  string
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
	if val, found := os.LookupEnv("SMTPHOST"); found {
		app.SMTP_HOST = val
		isRead = false
	}
	if val, found := os.LookupEnv("SMTPPORT"); found {
		cnv, _ := strconv.Atoi(val)
		app.SMTP_PORT = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("SMTPUSER"); found {
		app.SMTP_USER = val
		isRead = false
	}
	if val, found := os.LookupEnv("SMTPPASS"); found {
		app.SMTP_PASS = val
		isRead = false
	}
	if val, found := os.LookupEnv("PASSWDURL"); found {
		app.PASSWD_URL = val
		isRead = false
	}
	if val, found := os.LookupEnv("EMAILFROM"); found {
		app.EMAIL_FROM = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLIENTID"); found {
		CLIENT_ID = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLIENTSECRET"); found {
		CLIENT_SECRET = val
		isRead = false
	}
	if val, found := os.LookupEnv("GOOGLEURL"); found {
		GOOGLE_URL = val
		isRead = false
	}
	if val, found := os.LookupEnv("SCOPES"); found {
		SCOPES = strings.Split(val, ",")
		isRead = false
	}
	if val, found := os.LookupEnv("CLIENTIDFB"); found {
		CLIENT_ID_FB = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLIENTSECRETFB"); found {
		CLIENT_SECRET_FB = val
		isRead = false
	}
	if val, found := os.LookupEnv("FBURL"); found {
		FB_URL = val
		isRead = false
	}
	if val, found := os.LookupEnv("SCOPESFB"); found {
		SCOPES_FB = strings.Split(val, ",")
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
		SCOPES_FB = strings.Split(viper.GetString("SCOPESFB"), ",")
		FB_URL = viper.GetString("FBURL")
		CLIENT_ID_FB = viper.GetString("CLIENTIDFB")
		CLIENT_SECRET_FB = viper.GetString("CLIENTSECRETFB")
		SCOPES = strings.Split(viper.GetString("SCOPES"), ",")
		GOOGLE_URL = viper.GetString("GOOGLEURL")
		CLIENT_ID = viper.GetString("CLIENTID")
		CLIENT_SECRET = viper.GetString("CLIENTSECRET")
		AWS_ACCESS_KEY_ID = viper.GetString("AWSKEY")
		AWS_SECRET_ACCESS_KEY = viper.GetString("AWSSECRET")
		AWS_REGION = viper.GetString("AWSREGION")
		RDS_URL = viper.GetString("RDSURL")
		JWT_SECRET = viper.GetString("JWTSECRET")
		app.SMTP_HOST = viper.GetString("SMTPHOST")
		app.SMTP_PORT, _ = strconv.Atoi(viper.Get("SMTPPORT").(string))
		app.SMTP_USER = viper.GetString("SMTPUSER")
		app.SMTP_PASS = viper.GetString("SMTPPASS")
		app.PASSWD_URL = viper.GetString("PASSWDURL")
		app.EMAIL_FROM = viper.GetString("EMAILFROM")
		app.DB_USERNAME = viper.Get("DBUSER").(string)
		app.DB_PASSWORD = viper.Get("DBPASS").(string)
		app.DB_HOSTNAME = viper.Get("DBHOST").(string)
		app.DB_PORT, _ = strconv.Atoi(viper.Get("DBPORT").(string))
		app.DB_NAME = viper.Get("DBNAME").(string)
	}
	return &app
}
