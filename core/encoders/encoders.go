package encoders

import (
	"io"

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

type SavedTracksEncoder interface {
	Encode([]spotify.SavedTrack, io.Writer) error
}
