package player

import (
	"time"
	log "github.com/sirupsen/logrus"
	"io"
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

//go:embed beep.mp3
var fileBytes []byte
var player oto.Player
var decodedMp3 *mp3.Decoder
var otoCtx *oto.Context

func init() {
	fileBytesReader := bytes.NewReader(fileBytes)

	err := *new(error)
	decodedMp3, err = mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	samplingRate := 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo). 
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	readyChan := make(chan struct{})
	otoCtx, readyChan, err = oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

}

//play times with interval milliseconds in between
func Play(times, interval int) {
	log.Info("Playing")
	player = otoCtx.NewPlayer(decodedMp3)
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	count := 0;
	loop:
	for {
		select {
		case <-ticker.C:
			play(player)

			count++
			if count == times {
				ticker.Stop()
				break loop
			}
		}
	}
}

func play(player oto.Player) {
	_, err := player.(io.Seeker).Seek(0, io.SeekStart)
	if err != nil{
		panic("player.Seek failed: " + err.Error())
	}
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

func Close() {
	err := player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())

	}
}
