package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/sources"
)

func main() {
	// Set up configuration framework.
	viper.SetEnvPrefix("PRISMRIVER_")

	viper.SetDefault("DataDir", "/var/lib/prismriver")
	viper.SetDefault("Verbosity", "info")

	viper.BindEnv("DataDir")
	viper.BindEnv("Verbosity")

	verbosity := viper.GetString("Verbosity")
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		logrus.Errorf("Error reading verbosity level in configuration")
	}
	logrus.SetLevel(level)

	dataDir := viper.GetString("DataDir")
	os.MkdirAll(dataDir, os.ModeDir)

	playTest()
}

func playTest() {
	// My favorite Touhou Project arrangement (as of 2018/11/11).
	output, _ := sources.GetVideo("https://www.youtube.com/watch?v=24oZx-MTy68")
	file, err := os.Open(output)
	if err != nil {
		panic(err)
	}
	stream, format, err := vorbis.Decode(file)
	if err != nil {
		panic(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(stream, beep.Callback(func() {
		close(done)
	})))
	<-done
}
