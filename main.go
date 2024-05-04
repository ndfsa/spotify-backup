package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ndfsa/spotify-backup/auth"
	"github.com/ndfsa/spotify-backup/core"
	"github.com/zmb3/spotify/v2"
)

func main() {
	// get encoder from flags
	useAudioFeatures := flag.Bool("a", false, "include audio features in track dump")
	usePlaylist := flag.String("p", "", "dump playlist tracks, incompatible with '-f'")
	flag.Parse()

	ch := make(chan *spotify.Client)
	auth.SetupAuth(ch)
	spotifyClient := <-ch

	progressChannel := make(chan string)
	qChannel := make(chan int)
	go func() {
		for {
			var progress string
			select {
			case progress = <-progressChannel:
				fmt.Println(progress)
			case <-qChannel:
				return
			}
		}
	}()

	var tracks []core.DumpTrack
	if *usePlaylist != "" {
		playlist, err := core.GetPlaylist(spotifyClient, spotify.ID(*usePlaylist), progressChannel)
		if err != nil {
			log.Fatal(err)
		}

		tracks = playlist
	} else {
		favorites, err := core.GetFavorites(spotifyClient, progressChannel)
		if err != nil {
			log.Fatal(err)
		}

		tracks = favorites
	}

	currentDate := time.Now()
	fileName := fmt.Sprintf("trackdump-%d-%02d-%02d.json",
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day())

	if *useAudioFeatures {
		fullTracks, err := core.GetAudioFeatures(spotifyClient, tracks, progressChannel)
		if err != nil {
			log.Fatal(err)
		}
		core.WriteToFile(fullTracks, fileName)
	} else {
		core.WriteToFile(tracks, fileName)
	}
}
