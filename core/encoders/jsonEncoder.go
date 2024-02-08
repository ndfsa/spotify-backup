package encoders

import (
	"encoding/json"
	"io"

	"github.com/zmb3/spotify/v2"
)

type JsonEncoder struct {
}

func (enc JsonEncoder) Encode(items []spotify.SavedTrack, fields uint64, w io.Writer) error {
	tracks := []map[string]interface{}{}

	// extract relevant info
	for idx, entry := range items {
		current := map[string]interface{}{}
		if FieldNumber&fields != 0 {
			current[TagNumber] = idx + 1
		}
		if FieldAddedAt&fields != 0 {
			current[TagAddedAt] = entry.AddedAt
		}
		if FieldAlbum&fields != 0 {
			current[TagAlbum] = entry.Album.Name
		}
		if FieldArtists&fields != 0 {
			currentArtists := []string{}
			for _, artist := range entry.Artists {
				currentArtists = append(currentArtists, artist.Name)
			}
			current[TagArtists] = currentArtists
		}
		if FieldDuration&fields != 0 {
			current[TagDuration] = entry.Duration
		}
		if FieldExplicit&fields != 0 {
			current[TagExplicit] = entry.Explicit
		}
		if FieldId&fields != 0 {
			current[TagId] = entry.ID.String()
		}
		if FieldName&fields != 0 {
			current[TagName] = entry.Name
		}
		tracks = append(tracks, current)
	}

	jsonEnc := json.NewEncoder(w)
	if err := jsonEnc.Encode(&tracks); err != nil {
		return err
	}

	return nil
}
