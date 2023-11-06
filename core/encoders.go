package core

import (
	"encoding/csv"
	"io"
	"math/bits"
	"strconv"

	"github.com/zmb3/spotify/v2"
)

const (
	FieldAddedAt  uint = 1 << 0
	FieldAlbum    uint = 1 << 1
	FieldArtists  uint = 1 << 2
	FieldDuration uint = 1 << 3
	FieldExplicit uint = 1 << 4
	FieldId       uint = 1 << 5
	FieldName     uint = 1 << 6
)

type SavedTracksEncoder interface {
	Encode([]spotify.SavedTrack, uint, io.Writer) error
}

type CsvEncoder struct {
	Separator rune
}

func (enc CsvEncoder) Encode(items []spotify.SavedTrack, fields uint, w io.Writer) error {

	// create records array to encode later
	records := make([][]string, 0, len(items))

	// count number of fields each
	fieldCount := bits.OnesCount(fields)
	header := make([]string, 0, fieldCount)

	// create header dinamically
	if FieldAddedAt&fields != 0 {
		header = append(header, "added_at")
	}
	if FieldAlbum&fields != 0 {
		header = append(header, "album")
	}
	if FieldArtists&fields != 0 {
		header = append(header, "artist")
	}
	if FieldDuration&fields != 0 {
		header = append(header, "duration")
	}
	if FieldExplicit&fields != 0 {
		header = append(header, "explicit")
	}
	if FieldId&fields != 0 {
		header = append(header, "id")
	}
	if FieldName&fields != 0 {
		header = append(header, "name")
	}

	records = append(records, header)

	for _, item := range items {

		// whether or not it is needed to unwrap artist array
		unwrapArtists := false
		// keep track of where the artist field is located
		artistPos := 0

		// current record to build
		current := make([]string, 0, fieldCount)
		if FieldAddedAt&fields != 0 {
			current = append(current, item.AddedAt)
			artistPos++
		}
		if FieldAlbum&fields != 0 {
			current = append(current, item.Album.Name)
			artistPos++
		}
		if FieldArtists&fields != 0 {
			if len(item.Artists) == 1 {
				current = append(current, item.Artists[0].Name)
			} else {
				unwrapArtists = true
				current = append(current, "")
			}
		}
		if FieldDuration&fields != 0 {
			current = append(current, strconv.Itoa(item.Duration))
		}
		if FieldExplicit&fields != 0 {
			if item.Explicit {
				current = append(current, "1")
			} else {
				current = append(current, "0")
			}
		}
		if FieldId&fields != 0 {
			current = append(current, item.ID.String())
		}
		if FieldName&fields != 0 {
			current = append(current, item.Name)
		}

		if unwrapArtists {
			// repeat records for multiple artists
			for _, artist := range item.Artists {
				currentCopy := make([]string, fieldCount)
				copy(currentCopy, current)

				currentCopy[artistPos] = artist.Name
				records = append(records, currentCopy)
			}
		} else {
			records = append(records, current)
		}
	}

	// write to output stream
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = enc.Separator
	if err := csvWriter.WriteAll(records); err != nil {
		return err
	}
	csvWriter.Flush()

	return nil
}
