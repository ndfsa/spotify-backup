package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/zmb3/spotify/v2"
)

const (
	// track general info
	FieldNumber   uint64 = 1 << 0
	FieldAddedAt  uint64 = 1 << 1
	FieldAlbum    uint64 = 1 << 2
	FieldArtists  uint64 = 1 << 3
	FieldDuration uint64 = 1 << 4
	FieldExplicit uint64 = 1 << 5
	FieldId       uint64 = 1 << 6
	FieldName     uint64 = 1 << 7

	// track features

	// track analysis
)

type SavedTracksEncoder interface {
	Encode([]spotify.SavedTrack, uint64, io.Writer) error
}

type CsvEncoder struct {
	Separator rune
}

func (enc CsvEncoder) Encode(items []spotify.SavedTrack, fields uint64, w io.Writer) error {

	// create records array to encode later
	records := make([][]string, 0, len(items))

	// count number of fields each
	fieldCount := bits.OnesCount64(fields)
	header := make([]string, 0, fieldCount)

	// create header dinamically
	if FieldNumber&fields != 0 {
		header = append(header, "num")
	}
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

	for idx, item := range items {

		// whether or not it is needed to unwrap artist array
		unwrapArtists := false
		// keep track of where the artist field is located
		artistPos := 0

		// current record to build
		current := make([]string, 0, fieldCount)
		if FieldAddedAt&fields != 0 {
			current = append(current, fmt.Sprint(idx+1))
			artistPos++
		}
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
