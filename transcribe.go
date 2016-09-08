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
	"sync"
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
	//stoppedCh  chan bool
	//stop       bool
	wg sync.WaitGroup
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
	w.win = ui.NewWindow("Transcriber", 400, 240, false)
	w.picker = ui.NewEntry()
	w.picker.SetText("https://www.freesound.org/data/previews/" +
		"258/258397_450294-lq.mp3")
	w.slider = ui.NewSlider(0, 100)
	w.bStart = ui.NewButton("Start")
	w.bStart.OnClicked(func(*ui.Button) {
		// So pass true to the stop chan
		// When I want to end the UpdateSlider goroutine.
		w.stopCh <- true

		err := w.Start(w.picker.Text())
		if err != nil {
			fmt.Println(err)
		}
	})
	w.bPause = ui.NewButton("Pause")
	w.bPause.OnClicked(func(*ui.Button) {
		fmt.Println(w.TotalSeconds(), w.Seconds())
		if !w.transcribe.player.IsPlaying() {
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
	//w.stoppedCh = make(chan struct{})
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

func (mw *MainWindow) UpdateSlide() {
	for {
		seconds, err := mw.transcribe.player.MediaTime()
		if err != nil {
			fmt.Println(err)
		}
		if !(seconds > 0) {
			continue
		}
		fmt.Println(seconds)
		//seconds := millis * 1000
		//fmt.Println(seconds)
		break
	}
	//mw.slider = ui.NewSlider(0, seconds)
UpLoop:
	for {
		select {
		default:
			state, err := mw.transcribe.player.MediaState()
			if err != nil {
				fmt.Println("Get State Error: ", err)
				fmt.Println("Recording is not connected")
				break UpLoop
			}
			if state != 4 && state != 3 {
				continue UpLoop
			}
			pos, err := mw.transcribe.player.MediaPosition()
			if err != nil {
				fmt.Println(err)
				break UpLoop
			}
			percent := pos * 100
			mw.slider.SetValue(int(percent))
		case <-mw.stopCh:
			break UpLoop
		}
	}
	mw.wg.Done()
}

func (mw *MainWindow) Start(path string) error {
	// SetMedia for Player
	var err error
	// Don't 'start' until the goroutine updating GUI has stopped
	mw.wg.Wait()
	mw.wg.Add(1)
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

// Seconds returns an int of seconds total for media.
func (mw *MainWindow) TotalSeconds() int {
	seconds, err := mw.transcribe.player.MediaLength()
	if err != nil {
		fmt.Println(err)
	}
	return seconds / 1000
}

func (mw *MainWindow) Seconds() int {
	seconds, err := mw.transcribe.player.MediaTime()
	if err != nil {
		fmt.Println(err)
	}
	return seconds / 1000
}

// jumpBack jumps back in position.
// TODO: modify jump distance.
func (t *Transcriber) jumpBack() {
	pos, err := t.player.MediaTime()
	if err != nil {
		fmt.Println("Jump Back: ", err)
	}
	newPosition := pos - t.jump
	t.player.SetMediaTime(newPosition)
}

// jumpForward jumps forward position.
// TODO: modify jump distance.
func (t *Transcriber) jumpForward() {
	pos, err := t.player.MediaTime()
	if err != nil {
		fmt.Println("Jump Forward: ", err)
	}
	newPosition := pos + t.jump
	t.player.SetMediaTime(newPosition)
}
