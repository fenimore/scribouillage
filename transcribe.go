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
	bStart *ui.Button
	bPause *ui.Button
	bReset *ui.Button
	status *ui.Label
	box    *ui.Box
	slider *ui.Slider
	// Radio for jump value
	transcribe *Transcriber
	stopCh     chan bool
}

func NewTranscriber() *Transcriber {

	t := Transcriber{}
	// VLC data
	t.jump = 5000

	return &t
}

func NewMainWindow() *MainWindow {
	w := new(MainWindow)
	// User Interface
	w.win = ui.NewWindow("Transcriber", 400, 400, false)
	w.picker = ui.NewEntry()
	w.slider = ui.NewSlider(0, 100)
	w.bStart = ui.NewButton("Start")
	w.bStart.OnClicked(func(*ui.Button) {
		// Do nothign
		err := w.Start("https://www.freesound.org/data/previews/258/258397_450294-lq.mp3")
		if err != nil {
			fmt.Println(err)
		}
	})
	w.bPause = ui.NewButton("Pause")
	w.bPause.OnClicked(func(*ui.Button) {
		if w.transcribe.player.IsPlaying() {
			w.bPause.SetText("Pause")
		} else {
			w.bPause.SetText("Play")
		}
		w.transcribe.player.Pause(w.transcribe.player.IsPlaying())
	})
	w.status = ui.NewLabel("")
	w.box = ui.NewVerticalBox()
	w.box.Append(ui.NewLabel("Recording Path"), false)
	w.box.Append(w.picker, false)
	w.box.Append(w.slider, false)
	w.box.Append(w.bStart, false)
	w.box.Append(w.bPause, false)
	w.box.Append(w.status, false)
	w.win.SetChild(w.box)
	w.win.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	w.stopCh = make(chan bool)

	return w
}

func main() {
	t := NewTranscriber()
	err := vlc.Init("--no-video", "--quiet")
	if err != nil {
		fmt.Println("VLC init Error: %s\n", err)
	}
	defer vlc.Release()

	t.player, err = vlc.NewPlayer()
	if err != nil {
		fmt.Printf("VLC init Error: [%s]\nAre you using libvlc 2.x?\n", err)
		return
	}
	defer func() {
		t.player.Stop()
		t.player.Release()
	}()
	err = ui.Main(func() {
		mw := NewMainWindow()
		mw.transcribe = t
		mw.win.Show()
		err = mw.Start("/home/fen/flowers.mp3")
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (mw *MainWindow) UpdateSlide() error {
	for {
		state, err := mw.transcribe.player.GetState()
		if err != nil {
			fmt.Println("Get State Error: ", err)
			fmt.Println("Recording is not connected")
			return err
		}
		if state != 4 && state != 3 {
			continue
		}
		pos, err := mw.transcribe.player.GetPosition()
		if err != nil {
			return err
		}
		percent := pos * 100
		mw.slider.SetValue(int(percent))
	}
	return nil
}

func (mw *MainWindow) Start(path string) error {
	// SetMedia for Player
	var err error
	mw.transcribe.recording = path
	local := !strings.HasPrefix(mw.transcribe.recording, "http")
	if local {
		err = mw.transcribe.player.SetMedia(mw.transcribe.recording, true)
	} else {
		err = mw.transcribe.player.SetMedia(mw.transcribe.recording, false)
	}
	if err != nil {
		return err
	}
	err = mw.transcribe.player.Play()
	if err != nil {
		return err
	}
	// TODO: send a chan to stop last updateslide
	go mw.UpdateSlide()
	return nil
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
