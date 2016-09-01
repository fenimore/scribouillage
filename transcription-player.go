package main

import (
	"fmt"
	vlc "github.com/adrg/libvlc-go"
	"os"
	"os/exec"
)

func main() {
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println("File Open", err)
	}
	data := make([]byte, 24) // Buffer for reading file

	// VLC
	err = vlc.Init("--no-video", "--quiet")
	if err != nil {
		fmt.Println("Init", err)
	}
	// Defer defers in reverse order
	defer vlc.Release()

	player, err := vlc.NewPlayer()
	if err != nil {
		fmt.Println(err)
		return //give up
	}

	defer func() {
		player.Stop()
		player.Release()
	}()

	err = player.SetMedia("https://www.freesound.org/data/previews/258/258397_450294-lq.mp3", false)
	if err != nil {
		fmt.Println("Set Media", err)
		return
	}

	err = player.Play()
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
			fmt.Println("Left")
			prev()
			//jumpBack()
		case data[12]:
			fmt.Println("Center")
			err = player.Pause(player.IsPlaying())
			if err != nil {
				fmt.Println("Center", err)
			}
		case data[20]:
			fmt.Println("Right")
			next()
		}
	}
}

func next() {
	cmd := exec.Command("xte", "key XF86AudioNext")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func prev() {
	cmd := exec.Command("xte", "key XF86AudioPrev")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func pause() {
	cmd := exec.Command("xte", "key XF86AudioPlay")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// This is in the right track,
// But it totally doesn't work
// Shift/Alt/Shift+Alt small/medium/large jump
// according to default VLC
func jumpBack() {
	cmd := exec.Command("xte", "keydown Alt_L key Left keyup Alt_L")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func jumpForw() {
	cmd := exec.Command("xte", "keydown Alt_L key Right keyup Alt_L")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
