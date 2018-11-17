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
	viper.SetEnvPrefix("PRISMRIVER_")

	viper.SetDefault(constants.DATA, "/var/lib/prismriver")
	viper.SetDefault(constants.VERBOSITY, "info")

	viper.BindEnv(constants.DATA)
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
