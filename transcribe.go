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
	"fmt"
	"github.com/andlabs/ui"
	vlc "github.com/polypmer/libvlc-go"
	"strings"
)

type Transcriber struct {
	jump      int
	recording string
	player    *vlc.Player
}

type MainWindow struct {
	win    *ui.Window
	picker *ui.Entry
	bPlay  *ui.Button
	bPause *ui.Button
	bReset *ui.Button
	status *ui.Label
	box    *ui.Box
	slider *ui.Slider
	// Radio for jump value
}

func NewTranscriber() *Transcriber {
	t := new(Transcriber)

	// VLC data
	t.jump = 5000
	t.recording = "https://www.freesound.org/data/previews/258/258397_450294-lq.mp3"
	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		fmt.Println("VLC init Error: %s\n", err)
	}
	defer vlc.Release()
	t.player, err = vlc.NewPlayer()
	if err != nil {
		fmt.Printf("VLC init Error: [%s]\nAre you using libvlc 2.x?\n", err)
		return nil
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
		return nil
	}
	return t
}

func NewMainWindow() *MainWindow {
	w := new(MainWindow)
	// User Interface
	w.win = ui.NewWindow("Transcriber", 400, 400, false)
	w.picker = ui.NewEntry()
	w.slider = ui.NewSlider(0, 100)
	w.bPlay = ui.NewButton("Start")
	w.bPlay.OnClicked(func(*ui.Button) {
		// Do nothign
	})
	w.bPause = ui.NewButton("Pause")
	w.bPause.OnClicked(func(*ui.Button) {
		//t.player.Pause(t.player.IsPlaying())
	})
	w.status = ui.NewLabel("")
	w.box = ui.NewVerticalBox()
	w.box.Append(ui.NewLabel("Recording Path"), false)
	w.box.Append(w.picker, false)
	w.box.Append(w.slider, false)
	w.box.Append(w.bPlay, false)
	w.box.Append(w.bPause, false)
	w.box.Append(w.status, false)
	w.win.SetChild(w.box)
	w.win.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	return w
}

func main() {
	err := ui.Main(func() {
		mw := NewMainWindow()
		mw.win.Show()
		t := NewTranscriber()

		b := t.player.IsPlaying()
		fmt.Println(b)

	})
	if err != nil {
		fmt.Println(err)
	}
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
