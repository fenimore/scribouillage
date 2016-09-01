package main

import (
    "fmt"
    "time"

    vlc "github.com/adrg/libvlc-go"
)

func main() {
    // Initialize libvlc. Additional command line arguments can be passed in
    // to libvlc by specifying them in the Init function.
    if err := vlc.Init("--no-video", "--quiet"); err != nil {
        fmt.Println(err)
        return
    }
    defer vlc.Release()

    // Create a new player
    player, err := vlc.NewPlayer()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer func() {
        player.Stop()
        player.Release()
    }()

    // Set player media. The second parameter of the method specifies if
    // the media resource is local or remote.
    // err = player.SetMedia("localPath/test.mp4", true)
    err = player.SetMedia("http://stream-uk1.radioparadise.com/mp3-32", false)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Play
    err = player.Play()
    if err != nil {
        fmt.Println(err)
        return
    }

    time.Sleep(30 * time.Second)
}
