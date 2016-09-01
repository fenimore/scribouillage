package main

import (
	"fmt"
	vlc "github.com/polypmer/libvlc-go"
	"os"
)

type Transcriber struct {
	jump      int
	recording string
	player    *vlc.Player
}

func newTranscriber(path string) *Transcriber {
	return &Transcriber{
		jump:      5000,
		recording: path,
	}
}

func main() {
	// So-called driver for Infinity Pedal
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println("File Open", err)
	}
	data := make([]byte, 24) // Buffer for reading file

	// Transcriber Construction
	t := newTranscriber("https://www.freesound.org/" +
		"data/previews/258/258397_450294-lq.mp3")

	// VLC Init
	err = vlc.Init("--no-video", "--quiet")
	if err != nil {
		fmt.Println("Init", err)
	}
	defer vlc.Release()
	t.player, err = vlc.NewPlayer()
	if err != nil {
		fmt.Println(err)
		return //give up
	}
	//t.player = player
	defer func() {
		t.player.Stop()
		t.player.Release()
	}()
	err = t.player.SetMedia(t.recording, false)
	if err != nil {
		fmt.Println("Set Media", err)
		return
	}

	err = t.player.Play()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, err := file.Read(data)
		if err != nil {
			fmt.Println("Read", err)
		}
		switch byte(1) {
		case data[4]:
			t.jumpBack()
		case data[12]:
			fmt.Println("Center")
			err = t.player.Pause(t.player.IsPlaying())
			if err != nil {
				fmt.Println("Center", err)
			}
		case data[20]:
			t.jumpForward()
		}
	}
}

// jumpBack jumps back 5 seconds.
// TODO: modify jump distance.
func (t *Transcriber) jumpBack() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	newPosition := pos - 5000
	t.player.SetTime(newPosition)
}

// jumpForward jumps forward 5 seconds.
// TODO: modify jump distance.
func (t *Transcriber) jumpForward() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	newPosition := pos + 5000
	t.player.SetTime(newPosition)
}
