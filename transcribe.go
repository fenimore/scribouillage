package main

import (
	"fmt"
	ui "github.com/gizak/termui"
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
	// For debuggin
	t := newTranscriber("https://www.freesound.org/" +
		"data/previews/258/258397_450294-lq.mp3")

	// VLC Initialization and Player Construction
	err = vlc.Init("--no-video", "--quiet")
	if err != nil {
		fmt.Println("VLC init Error: %s\n", err)
	}
	defer vlc.Release()
	t.player, err = vlc.NewPlayer()
	if err != nil {
		fmt.Println(err)
		return //give up
	}
	defer func() {
		t.player.Stop()
		t.player.Release()
	}()
	err = t.player.SetMedia(t.recording, false)
	if err != nil {
		fmt.Println("Set Media Error: %s\n", err)
		return
	}

	// Start Playing Recording
	err = t.player.Play()
	if err != nil {
		fmt.Println(err)
		return
	}

	// User Interface
	fmt.Println(t.stats())
	fmt.Println("Start Recording: ")
	err = ui.Init()
	if err != nil {
		fmt.Println(err)
	}
	defer ui.Close()

	g := ui.NewGauge()
	g.Percent = 0
	g.Width = 50
	g.Height = 3
	//g.Y = 14
	g.BorderLabel = "Recording"
	g.Label = "{{percent}}%"

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Render(g)
	go func() {
		for {
			pos, err := t.player.GetPosition()
			if err != nil {
				fmt.Println(err)
			}
			percentage := pos * 100
			g.Percent = int(percentage)
			ui.Render(g)
		}
	}()
	go func() {
		for {
			_, err := file.Read(data)
			if err != nil {
				fmt.Println("Reading Hiddev Error: %s\n", err)
			}
			switch byte(1) {
			case data[4]:
				t.jumpBack()
			case data[12]:
				err = t.player.Pause(t.player.IsPlaying())
				if err != nil {
					fmt.Printf("Center Error: %s\n", err)
				}
				fmt.Println(t.stats())
			case data[20]:
				t.jumpForward()
			}
		}
	}()
	ui.Loop()
}

// jumpBack jumps back in position.
// TODO: modify jump distance.
func (t *Transcriber) jumpBack() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	newPosition := pos - t.jump
	t.player.SetTime(newPosition)
}

// jumpForward jumps forward position.
// TODO: modify jump distance.
func (t *Transcriber) jumpForward() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	newPosition := pos + t.jump
	t.player.SetTime(newPosition)
}

func (t *Transcriber) stats() string {
	tim, err := t.player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	pos, err := t.player.GetPosition()
	if err != nil {
		fmt.Println(err)
	}
	len, err := t.player.GetLength()
	if err != nil {
		fmt.Println(err)
	}
	second := tim / 1000
	percentage := pos * 100
	total := len / 1000
	return fmt.Sprintf("\nOf %d seconds, in second %d\nPercentage: %.f%%",
		total, second, percentage)
}
