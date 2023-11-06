package main

import (
	"log"

	"github.com/ndfsa/spotify-backup/auth"
	"github.com/ndfsa/spotify-backup/core"
	"github.com/zmb3/spotify/v2"
)

func main() {
	ch := make(chan *spotify.Client)
	auth.SetupAuth(ch)

	// wait for client creation
	client := <-ch

	favorites, err := core.GetFavorites(client)
	if err != nil {
		log.Fatal(err)
	}

	core.WriteToFile(
		favorites,
		core.FieldAddedAt|
			core.FieldAlbum|
			core.FieldArtists|
			core.FieldDuration|
			core.FieldExplicit|
			core.FieldId|
			core.FieldName,
		core.CsvEncoder{Separator: '\t'},
		"backup.csv")
}
