package encoders

import (
	"encoding/json"
	"io"

	"github.com/zmb3/spotify/v2"
)

type JsonEncoder struct {
}

func (enc JsonEncoder) Encode(items []spotify.SavedTrack, w io.Writer) error {
	tracks := []map[string]interface{}{}

	// extract relevant info
	for idx, entry := range items {
		current := map[string]interface{}{}
		current[TagNumber] = idx + 1
		current[TagAddedAt] = entry.AddedAt
		current[TagAlbum] = entry.Album.Name
		currentArtists := []string{}
		for _, artist := range entry.Artists {
			currentArtists = append(currentArtists, artist.Name)
		}
		current[TagArtists] = currentArtists
		current[TagDuration] = entry.Duration
		current[TagExplicit] = entry.Explicit
		current[TagId] = entry.ID.String()
		current[TagName] = entry.Name
		tracks = append(tracks, current)
	}

	jsonEnc := json.NewEncoder(w)
	if err := jsonEnc.Encode(&tracks); err != nil {
		return err
	}

	return nil
}
