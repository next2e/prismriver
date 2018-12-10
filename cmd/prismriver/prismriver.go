package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/server"
	"os"
)

func main() {
	// Set up configuration framework.
	viper.SetEnvPrefix("PRISMRIVER")
	viper.AutomaticEnv()

	viper.SetDefault(constants.DATA, "/var/lib/prismriver")
	viper.SetDefault(constants.DBHOST, "localhost")
	viper.SetDefault(constants.DBNAME, "prismriver")
	viper.SetDefault(constants.DBPASSWORD, "prismriver")
	viper.SetDefault(constants.DBPORT, "5432")
	viper.SetDefault(constants.DBUSER, "prismriver")
	viper.SetDefault(constants.VERBOSITY, "info")

	viper.BindEnv(constants.DATA)
	viper.BindEnv(constants.DBHOST)
	viper.BindEnv(constants.DBNAME)
	viper.BindEnv(constants.DBPASSWORD)
	viper.BindEnv(constants.DBPORT)
	viper.BindEnv(constants.DBUSER)
	viper.BindEnv(constants.VERBOSITY)

	verbosity := viper.GetString(constants.VERBOSITY)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		logrus.Errorf("Error reading verbosity level in configuration")
	}
	logrus.SetLevel(level)

	dataDir := viper.GetString(constants.DATA)
	os.MkdirAll(dataDir, os.ModeDir)

	server.CreateRouter()
}
