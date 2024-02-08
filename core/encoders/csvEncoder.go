package encoders

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/zmb3/spotify/v2"
)

const (
	TagNumber   string = "num"
	TagAddedAt  string = "added_at"
	TagAlbum    string = "album"
	TagArtists  string = "artists"
	TagDuration string = "duration"
	TagExplicit string = "explicit"
	TagId       string = "id"
	TagName     string = "name"
)

type CsvEncoder struct {
	Separator rune
}

func (enc CsvEncoder) Encode(items []spotify.SavedTrack, fields uint64, w io.Writer) error {

	// create records array to encode later
	records := make([][]string, 0, len(items))

	// count number of fields
	fieldCount := bits.OnesCount64(fields)
	header := make([]string, 0, fieldCount)

	// create header
	if FieldNumber&fields != 0 {
		header = append(header, TagNumber)
	}
	if FieldAddedAt&fields != 0 {
		header = append(header, TagAddedAt)
	}
	if FieldAlbum&fields != 0 {
		header = append(header, TagAlbum)
	}
	if FieldArtists&fields != 0 {
		header = append(header, TagArtists)
	}
	if FieldDuration&fields != 0 {
		header = append(header, TagDuration)
	}
	if FieldExplicit&fields != 0 {
		header = append(header, TagExplicit)
	}
	if FieldId&fields != 0 {
		header = append(header, TagId)
	}
	if FieldName&fields != 0 {
		header = append(header, TagName)
	}

	records = append(records, header)

	for idx, item := range items {

		// whether or to unwrap the artists list into separate entries
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
			// repeat records when the track has multiple artists
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

	return nil
}
