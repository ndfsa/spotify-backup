package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ndfsa/spotify-backup/auth"
	"github.com/ndfsa/spotify-backup/core"
	"github.com/ndfsa/spotify-backup/core/encoding"
	"github.com/zmb3/spotify/v2"
)

func main() {
	// get encoder from flags
	useAudioFeatures := flag.Bool("A", false, "include audio features in backup")
	flag.Parse()

	ch := make(chan *spotify.Client)
	auth.SetupAuth(ch)

	client := <-ch

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

	favorites, err := core.GetFavorites(client, progressChannel)
	if err != nil {
		log.Fatal(err)
	}

	tracks := make([]map[string]interface{}, 0, len(favorites))
	encoding.EncodeSavedTracks(favorites, &tracks)

	if *useAudioFeatures {
		audioFeatures, err := core.GetAudioFeatures(client, favorites, progressChannel)
		if err != nil {
			log.Fatal(err)
		}
		encoding.EncodeAudioFeatures(audioFeatures, &tracks)

	}

	currentDate := time.Now()
	fileName := fmt.Sprintf("backup-%d-%02d-%02d.json",
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day())
	core.WriteToFile(tracks, fileName)
}
