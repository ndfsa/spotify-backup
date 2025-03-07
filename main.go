package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"path"
	"time"

	"spotify-backup/auth"
	"spotify-backup/core"

	"github.com/zmb3/spotify/v2"
)

var versionId = "dev"

func main() {
	// get encoder from flags
	useAudioFeatures := flag.Bool("a", false, "include audio features in track dump")
	usePlaylist := flag.String("p", "", "dump playlist tracks")
	version := flag.Bool("v", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println(versionId)
		return
	}

	spotifyClient := auth.GetClient()

	progressChannel := make(chan int)
	go func() {
		for {
			fmt.Printf("progress: %d\n", <-progressChannel)
		}
	}()

	var tracks []core.DumpTrack
	var prefix string
	if *usePlaylist != "" {
		uri, err := url.Parse(*usePlaylist)

		var playlistId string
		if err != nil {
			playlistId = *usePlaylist
		} else {
			playlistId = path.Base(uri.Path)
		}

		fmt.Println("Dumping playlist tracks")
		playlist, name, err := core.GetPlaylist(
			spotifyClient,
			spotify.ID(playlistId),
			progressChannel)
		if err != nil {
			log.Fatal(err)
		}

		tracks = playlist
		prefix = fmt.Sprintf("playlist-(%s)", name)
	} else {
		fmt.Println("Dumping favorite tracks")
		favorites, err := core.GetFavorites(spotifyClient, progressChannel)
		if err != nil {
			log.Fatal(err)
		}

		tracks = favorites
		prefix = "favorites"
	}

	currentDate := time.Now()
	fileName := fmt.Sprintf("%d-%02d-%02d.json",
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day())

	if *useAudioFeatures {
		fmt.Println("Adding audio features to dump data")
		fullTracks, err := core.GetAudioFeatures(spotifyClient, tracks, progressChannel)
		if err != nil {
			log.Fatal(err)
		}
		core.WriteToFile(fullTracks, fmt.Sprintf("%s-full-%s", prefix, fileName))
	} else {
		core.WriteToFile(tracks, fmt.Sprintf("%s-%s", prefix, fileName))
	}
}
