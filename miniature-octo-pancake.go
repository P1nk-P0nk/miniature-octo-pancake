package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

func main() {

	var filename string
	var streamer beep.StreamSeekCloser
	var format beep.Format

	flag.StringVar(&filename, "filename", "", "The path to the filename you wish to open")

	flag.Parse()

	fmt.Println(filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	mime, err := mimetype.DetectReader(f)
	switch mime.Extension() {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	case ".wav":
		streamer, format, err = wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
	case ".flac":
		streamer, format, err = flac.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			speaker.Lock()
			fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
			speaker.Unlock()
		}
	}

}
