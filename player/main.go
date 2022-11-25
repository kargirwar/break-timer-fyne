package player

import (
	"time"
	"log"
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/tevino/abool/v2"
)

//go:embed beep.mp3
var fileBytes []byte
var isPlaying *abool.AtomicBool
var player oto.Player

func init() {
	isPlaying = abool.New()
	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
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
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	player = otoCtx.NewPlayer(decodedMp3)
}

func Play() {
	if isPlaying.IsSet() {
		log.Println("Already playing")
		return
	}

	go play()

	for isPlaying.IsNotSet() {
		time.Sleep(time.Millisecond)
	}
}

func play() {
	log.Println("Playing");
	isPlaying.Set()
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	isPlaying.UnSet();
	log.Println("Done");
}

func Close() {
	for isPlaying.IsSet() {
		time.Sleep(time.Millisecond)
	}

	log.Println("Closing");
	err := player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())

	}
}
