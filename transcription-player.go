package main

import (
	"fmt"
	vlc "github.com/adrg/libvlc-go"
	"os"
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
		case data[12]:
			fmt.Println("Center")
			err = player.Pause(player.IsPlaying())
			if err != nil {
				fmt.Println("Center", err)
			}
			fmt.Println(position(player))
			fmt.Println(time(player))
			fmt.Println(length(player))
		case data[20]:
			fmt.Println("Right")
		}
	}
}

func length(player *vlc.Player) int {
	duration, err := player.GetLength()
	if err != nil {
		fmt.Println(err)
	}
	return duration
}

func time(player *vlc.Player) int {
	t, err := player.GetTime()
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func position(player *vlc.Player) float32 {
	p, err := player.GetPosition()
	if err != nil {
		fmt.Println(err)
	}
	return p
}
