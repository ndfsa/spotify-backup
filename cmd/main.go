package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ndfsa/spotify-backup/auth"
	"github.com/ndfsa/spotify-backup/core"
	"github.com/zmb3/spotify/v2"
)

func main() {
	ch := make(chan *spotify.Client)
	auth.SetupAuth(ch)

	client := <-ch

	favorites, err := core.GetFavorites(client)
	if err != nil {
		log.Fatal(err)
	}

	currentDate := time.Now()
	fileName := fmt.Sprintf("backup-%d-%02d-%02d.csv",
		currentDate.Year(),
		currentDate.Month(),
		currentDate.Day())

	core.WriteToFile(
		favorites,
		core.FieldNumber|
			core.FieldAddedAt|
			core.FieldAlbum|
			core.FieldArtists|
			core.FieldDuration|
			core.FieldExplicit|
			core.FieldId|
			core.FieldName,
		core.CsvEncoder{Separator: '\t'},
		fileName)
}
