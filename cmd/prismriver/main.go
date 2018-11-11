package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/ttpcodes/prismriver/internal/app/sources"
)

func main() {
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
