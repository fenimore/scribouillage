// package transcribe uses libvlc bindings to play a recording.
// The driver, so-called, then communicates with libvlc in order
// to jump forward and back in the recording.
//
// Fenimore Love (c) 2016 | GPLv3
//
// TODO:
// - add flags for jump values.
// - finish terminal ui.
package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui"
	vlc "github.com/polypmer/libvlc-go"
	"os"
	"strings"
	"time"
)

type Transcriber struct {
	jump      int
	recording string
	player    *vlc.Player
}

func newTranscriber(path string, seconds int) *Transcriber {
	seconds = seconds * 1000
	return &Transcriber{
		jump:      seconds,
		recording: path,
	}
}

func main() {
	// So-called driver for Infinity Pedal
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println("Footpedal Error: ", err)
		fmt.Println("Make sure the Footpedal is plugged in to your computer.\nhidddev0 should appear in /dev/usb, and you must have reading privilegs for this file.")
		return
	}
	data := make([]byte, 24) // Buffer for reading file

	pathFlag := flag.String("p", "https://www.freesound.org/data/previews/258/258397_450294-lq.mp3", "Path to recording.")
	jumpFlag := flag.Int("j", 5, "Jump distance in seconds.")
	flag.Parse()
	// Transcriber Construction
	// For debuggin
	// TODO: local files don't work?
	t := newTranscriber(*pathFlag, *jumpFlag)
	fmt.Println(t.jump, *jumpFlag)
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
	// SetMedia for Player
	local := !strings.HasPrefix(t.recording, "http")
	if local {
		err = t.player.SetMedia(t.recording, true)
	} else {
		err = t.player.SetMedia(t.recording, false)
	}
	if err != nil {
		fmt.Println("Set Media Path Error: %s\n", err)
		return
	}
	// Start Playing Recording
	err = t.player.Play()
	if err != nil {
		fmt.Println("Player Play: ", err)
		return
	}

	// User Interface
	fmt.Println("Start Recording: ")

	// UI thread

	go t.driverThread(data, file)
	t.uiThread()
	//ui.Loop()
}

// jumpBack jumps back in position.
// TODO: modify jump distance.
func (t *Transcriber) jumpBack() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println("Jump Back: ", err)
	}
	newPosition := pos - t.jump
	t.player.SetTime(newPosition)
}

// jumpForward jumps forward position.
// TODO: modify jump distance.
func (t *Transcriber) jumpForward() {
	pos, err := t.player.GetTime()
	if err != nil {
		fmt.Println("Jump Forward: ", err)
	}
	newPosition := pos + t.jump
	t.player.SetTime(newPosition)
}

func (t *Transcriber) stats() (string, string) {
	tim, err := t.player.GetTime()
	if err != nil {
		fmt.Println("Stats Get Time: ", err)
	}
	len, err := t.player.GetLength()
	if err != nil {
		fmt.Println("Stats Get Length: ", err)
	}
	second := tim / 1000
	total := len / 1000

	return fmt.Sprintf("%d", second), fmt.Sprintf("%d", total)
}

// Run as goroutine?
func (t *Transcriber) driverThread(buff []byte, file *os.File) {
	for {
		_, err := file.Read(buff)
		if err != nil {
			fmt.Println("Reading Hiddev Error: %s\n", err)
			fmt.Println("You must have reading privilegs for the file hiddev0 in /dev/usb/")
			return
		}
		switch byte(1) {
		case buff[4]:
			t.jumpBack()
		case buff[12]:
			//fmt.Println(t.player.IsPlaying())
			err = t.player.Pause(t.player.IsPlaying())
			if err != nil {
				fmt.Printf("Center Error: %s\n", err)
			}
			//fmt.Println(t.stats())
		case buff[20]:
			t.jumpForward()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (t *Transcriber) uiThread() {
	err := ui.Init()
	if err != nil {
		fmt.Println("UI Init: ", err)
	}
	defer ui.Close()

	// Recording Progress
	g := ui.NewGauge()
	g.Percent = 0
	g.Width = 50
	g.Height = 3
	//g.Y = 14
	g.BorderLabel = "Recording"
	g.Label = "{{percent}}%"

	// Instructions
	p := ui.NewPar("Press q to quit")
	p.Height = 3
	p.Width = 50
	p.Y = 6
	p.BorderLabel = "Instructions"
	p.BorderFg = ui.ColorYellow

	// Status
	s := ui.NewPar(" /  seconds")
	s.Height = 3
	s.Width = 50
	s.Y = 3
	s.BorderLabel = "Time"

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Render(g, p, s)
	go func() {
		for {
			state, err := t.player.GetState()
			if err != nil {
				fmt.Println("Get State: ", err)
				break
			}
			if state != 4 && state != 3 && state != 6 {
				continue
			}
			pos, err := t.player.GetPosition()
			if err != nil {
				fmt.Println("Get Position: ", err)
			}
			percentage := pos * 100
			g.Percent = int(percentage)
			sec, tot := t.stats()
			status := fmt.Sprintf("%s / %s seconds", sec, tot)
			s = ui.NewPar(status)
			s.Height = 3
			s.Width = 50
			s.Y = 3
			s.BorderLabel = "Time"
			ui.Render(g, s)
			if state == 6 {
				err = t.player.Stop()
				if err != nil {
					fmt.Println("Stop: ", err)
				}
			}
			//fmt.Println("Blip")
		}
		fmt.Println("No active Connection?")
	}()
	ui.Loop()
}
