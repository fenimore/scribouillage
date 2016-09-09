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
	"encoding/binary"
	"fmt"
	"github.com/andlabs/ui"
	vlc "github.com/polypmer/libvlc-go"
	"github.com/zserge/hid"
	"strings"
	"sync"
	"time"
)

const (
	left    = 1
	right   = 4
	middle  = 2
	release = 0
)

var pedalId string = "05f3:00ff" // Vendor and Product for Infinity Pedal

// The transcriber, part of the MainWindow
// keeps track of the recording to be transcribed.
type Transcriber struct {
	jump      int
	recording string
	player    *vlc.Player
}

// MainWindow is the main GUI window
type MainWindow struct {
	// TODO add number picker or radiobox for jump.
	win      *ui.Window
	picker   *ui.Entry
	bStart   *ui.Button
	bPause   *ui.Button
	bReset   *ui.Button
	bForw    *ui.Button
	bBack    *ui.Button
	lTotal   *ui.Label
	lCurrent *ui.Label
	box      *ui.Box
	controls *ui.Box
	location *ui.Box
	seeks    *ui.Box
	slider   *ui.Slider
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
	w.win = ui.NewWindow("Scribouillage", 500, 200, true)
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
		if !w.transcribe.player.IsPlaying() {
			w.bPause.SetText("Pause")
		} else {
			w.bPause.SetText("Play")
		}
		w.transcribe.player.Pause(
			w.transcribe.player.IsPlaying())
	})
	w.bForw = ui.NewButton("Forward >")
	w.bForw.OnClicked(func(*ui.Button) {
		w.transcribe.jumpForward()
		t, err := w.transcribe.player.MediaTime()
		if err != nil {
			fmt.Println(err)
		}
		w.lCurrent.SetText(w.Minutes(t))
	})
	w.bBack = ui.NewButton("< Back")
	w.bBack.OnClicked(func(*ui.Button) {
		w.transcribe.jumpBack()
		t, err := w.transcribe.player.MediaTime()
		if err != nil {
			fmt.Println(err)
		}
		w.lCurrent.SetText(w.Minutes(t))
	})
	w.lTotal = ui.NewLabel("")
	w.lCurrent = ui.NewLabel("")
	// Boxes goodies
	w.box = ui.NewVerticalBox()
	w.controls = ui.NewHorizontalBox()
	w.location = ui.NewHorizontalBox()
	w.seeks = ui.NewHorizontalBox()
	// File Picker
	w.box.Append(ui.NewLabel("Path to recording:"), false)
	w.controls.Append(w.picker, true)
	w.controls.Append(w.bStart, true)
	w.box.Append(w.controls, false)
	//w.box.Append(ui.NewHorizontalSeparator(), false)//dontwork
	// Location and slider
	w.location.Append(w.lCurrent, false)
	w.location.Append(w.lTotal, false)
	w.box.Append(w.location, false)
	w.box.Append(w.slider, false)
	// Start and Play controls
	w.box.Append(w.bPause, false)
	// Seek navigation
	w.seeks.Append(w.bBack, true)
	w.seeks.Append(w.bForw, true)
	w.box.Append(w.seeks, false)
	// Set box to window
	w.win.SetChild(w.box)
	//w.win.SetChild(w.controls)

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
		go mw.LaunchDriver()
		// TODO: Start Foot Pedal Driver here?
		// go footpedal?
	})
	if err != nil {
		fmt.Println(err)
	}
}

// UpdateSlide, run as goroutine, updates GUI slide.
// To cancel, pass true into MainWindow.stopCh chan.
func (mw *MainWindow) UpdateSlide() {
	for {
		//		plac, err := mw.transcribe.player.MediaTime()
		//		if err != nil {
		//			fmt.Println(err)
		//		}
		leng, err := mw.transcribe.player.MediaLength()
		if err != nil {
			fmt.Println(err)
		}
		if !(leng > 0) {
			continue
		}
		//mw.lCurrent.SetText(mw.Minutes(plac))
		mw.lTotal.SetText(" " + mw.Minutes(leng))
		break
	}
	// set length to slider scale
	//mw.slider = ui.NewSlider(0, length)
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
			if state != 4 && state != 3 && state != 6 {
				continue UpLoop
			}
			pos, err := mw.transcribe.player.MediaPosition()
			if err != nil {
				fmt.Println(err)
				break UpLoop
			}
			//t, err := mw.transcribe.player.MediaTime()
			//if err != nil {
			//	fmt.Println(err)
			//break UpLoop
			//}// This breaks it?!??!
			//fmt.Println(mw.Minutes(t))
			//mw.lCurrent.SetText(mw.Minutes(t))
			percent := pos * 100
			mw.slider.SetValue(int(percent))
		case <-mw.stopCh:
			break UpLoop
		}
	}
	mw.wg.Done()
}

// Start sets media to path and plays recording.
// There is a sync lock because this method calls
// UpdateSlide, and only one of these goroutines
// should be running at a time.
func (mw *MainWindow) Start(path string) error {
	// SetMedia for Player
	var err error
	// Don't 'start' until the goroutine updating GUI has stopped
	mw.wg.Wait()
	mw.wg.Add(1)
	mw.transcribe.recording = path
	local := !strings.HasPrefix(mw.transcribe.recording, "http")
	if local {
		err = mw.transcribe.player.SetMedia(
			mw.transcribe.recording, true)
	} else {
		err = mw.transcribe.player.SetMedia(
			mw.transcribe.recording, false)
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

// Converts milliseconds to Minutes and returns string.
func (mw *MainWindow) Minutes(length int) string {
	s := length / 1000
	return fmt.Sprintf("%d:%d", s/60, s%60)
}

// jumpBack jumps back in position.
// TODO: modify jump distance.
// TODO: change name to Backward?
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

// Find Pedal Device and Process Pedal Presses.
// Launch in Goroutine, obs.
func (mw *MainWindow) LaunchDriver() {
	var device hid.Device
	hid.UsbWalk(func(dev hid.Device) {
		info := dev.Info()
		d := fmt.Sprintf("%04x:%04x", info.Vendor, info.Product)
		if d == pedalId {
			device = dev
		}
	})
	if device == nil {
		fmt.Println("Error: no device found.")
		return
	}
	err := device.Open()
	if err != nil {
		fmt.Printf("Open Error: %s\nCheck Privileges\n", err)
		return
	}
	defer device.Close()
	for {
		buf, err := device.Read(-1, 1*time.Second)
		if err == nil {
			// otherwise get err 'connection timed out'
			switch binary.LittleEndian.Uint16(buf) {
			case left:
				mw.transcribe.jumpBack()
			case right:
				mw.transcribe.jumpForward()
			case middle:
				mw.transcribe.player.Pause(
					mw.transcribe.player.IsPlaying())
			case release:
				// do nothing
				// Unless only play when pressed
				// is set.
				continue
			default:
				// do nothing
				continue
			}
		}
	}
}
