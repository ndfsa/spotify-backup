package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ndfsa/spotify-backup/auth"
	"github.com/ndfsa/spotify-backup/core"
	"github.com/ndfsa/spotify-backup/core/encoders"
	"github.com/zmb3/spotify/v2"
)

func main() {
	// get encoder from flags
	encoding := flag.String("encoder", "json", "encoder to use, only json and csv supported")
	flag.Parse()

	var encoder encoders.SavedTracksEncoder

	switch *encoding {
	case "json":
		encoder = encoders.JsonEncoder{}
	case "csv":
		encoder = encoders.CsvEncoder{Separator: '\t'}
	default:
		log.Fatal(errors.New("unknown encoder: " + *encoding))
	}
	// get extension for file
	ext := *encoding

	ch := make(chan *spotify.Client)
	auth.SetupAuth(ch)

	client := <-ch

	favorites, err := core.GetFavorites(client)
	if err != nil {
		log.Fatal(err)
	}

	currentDate := time.Now()
	fileName := fmt.Sprintf("backup-%d-%02d-%02d.%s",
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day(),
		ext)

	core.WriteToFile(
		favorites,
		encoders.FieldNumber|
			encoders.FieldAddedAt|
			encoders.FieldAlbum|
			encoders.FieldArtists|
			encoders.FieldDuration|
			encoders.FieldExplicit|
			encoders.FieldId|
			encoders.FieldName,
		encoder,
		fileName)
}
