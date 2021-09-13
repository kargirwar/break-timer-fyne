package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var buffer *beep.Buffer

func init() {
	f, err := os.Open("beep.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/60))

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

}

func play() {
	log.Println("Starting play thread")
	//todo: is there a way to create a stopped ticker?
	ticker := time.NewTicker(time.Second)
	ticker.Stop()

	for {
		select {
		case cmd := <-playerCh:
			if cmd == PLAY {
				log.Println("Playing")
				ticker = time.NewTicker(time.Second)
			} else if cmd == STOP {
				log.Println("Stopping")
				ticker.Stop()
			}
		case <-ticker.C:
			shot := buffer.Streamer(0, buffer.Len())
			speaker.Play(shot)
		}
	}
}
