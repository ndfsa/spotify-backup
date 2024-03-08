package encoders

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/zmb3/spotify/v2"
)

type CsvEncoder struct {
	Separator rune
}

func (enc CsvEncoder) Encode(items []spotify.SavedTrack, w io.Writer) error {

	// create records array to encode later
	records := make([][]string, 0, len(items))

	// create header
	header := []string{TagNumber,
		TagAddedAt,
		TagAlbum,
		TagArtists,
		TagDuration,
		TagExplicit,
		TagId,
		TagName}
	fieldCount := len(header)

	records = append(records, header)

	for idx, item := range items {
		record := make([]string, fieldCount)

		record = append(record,
			fmt.Sprint(idx+1),
			item.AddedAt,
			parseArtists(item.Artists),
			fmt.Sprint(item.Duration),
			fmt.Sprint(item.Explicit),
			item.ID.String(),
			item.Name)
	}

	// write to output stream
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = enc.Separator
	if err := csvWriter.WriteAll(records); err != nil {
		return err
	}

	return nil
}

func parseArtists(artists []spotify.SimpleArtist) string {
	length := len(artists)
	if length < 1 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(artists[0].Name)
	if length > 1 {
		for _, artist := range artists[1:] {
			builder.WriteString(", ")
			builder.WriteString(artist.Name)
		}
	}

	return builder.String()
}
