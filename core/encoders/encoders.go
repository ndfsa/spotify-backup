package encoders

import (
	"io"

	"github.com/zmb3/spotify/v2"
)

const (
	// track general info
	FieldNumber uint64 = 1 << iota
	FieldAddedAt
	FieldAlbum
	FieldArtists
	FieldDuration
	FieldExplicit
	FieldId
	FieldName

	// track features

	// track analysis
)

type SavedTracksEncoder interface {
	Encode([]spotify.SavedTrack, uint64, io.Writer) error
}
