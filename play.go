package main

import (
	"os"
	"time"
	"context"
	"fmt"

	"log"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/kargirwar/golang/utils"
	"github.com/kargirwar/golang/utils/macos"
)

var buffer *beep.Buffer

func init() {
	f, err := os.Open("/usr/local/bin/beep.mp3")
	if err != nil {
		utils.Dbg(context.Background(), err.Error())
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		utils.Dbg(context.Background(), err.Error())
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/60))

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

}

func play() {
	utils.Dbg(context.Background(), fmt.Sprintln("Starting play thread"))
	//todo: is there a way to create a stopped ticker?
	ticker := time.NewTicker(time.Second)
	ticker.Stop()

	for {
		select {
		case cmd := <-playerCh:
			if cmd == PLAY {
				utils.Dbg(context.Background(), fmt.Sprintln("Playing"))
				ticker = time.NewTicker(time.Second)
			} else if cmd == STOP {
				utils.Dbg(context.Background(), fmt.Sprintln("Stopping"))
				ticker.Stop()
			}
		case <-ticker.C:
			locked, err := macos.IsScreenLocked(context.Background())
			if err == nil {
				if !locked {
					shot := buffer.Streamer(0, buffer.Len())
					speaker.Play(shot)
					continue
				}

				utils.Dbg(context.Background(), "Screen locked")
			} else {
				utils.Dbg(context.Background(), err.Error())
			}
		}
	}
}
