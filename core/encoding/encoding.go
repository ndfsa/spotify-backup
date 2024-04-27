package encoding

import (
	"github.com/zmb3/spotify/v2"
)

const (
	TagAddedAt  string = "added_at"
	TagAlbum    string = "album"
	TagArtists  string = "artists"
	TagDuration string = "duration"
	TagExplicit string = "explicit"
	TagId       string = "id"
	TagName     string = "name"

	TagAcousticness     string = "acousticness"
	TagDanceability     string = "danceability"
	TagEnergy           string = "energy"
	TagInstrumentalness string = "instrumentalness"
	TagKey              string = "key"
	TagLiveness         string = "liveness"
	TagLoudness         string = "loudness"
	TagMode             string = "mode"
	TagSpeechiness      string = "speechiness"
	TagTempo            string = "tempo"
	TagTimeSignature    string = "time_signature"
	TagValence          string = "valence"
)

func EncodeSavedTracks(items []spotify.SavedTrack, tracks *[]map[string]interface{}) {
	// extract relevant info
	for _, entry := range items {
		current := map[string]interface{}{}
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
		*tracks = append(*tracks, current)
	}
}

func EncodeAudioFeatures(items []*spotify.AudioFeatures, tracks *[]map[string]interface{}) {
	// extract only relevant info
	for idx, entry := range items {
		(*tracks)[idx][TagAcousticness] = entry.Acousticness
		(*tracks)[idx][TagDanceability] = entry.Danceability
		(*tracks)[idx][TagEnergy] = entry.Energy
		(*tracks)[idx][TagInstrumentalness] = entry.Instrumentalness
		(*tracks)[idx][TagKey] = entry.Key
		(*tracks)[idx][TagLiveness] = entry.Liveness
		(*tracks)[idx][TagLoudness] = entry.Loudness
		(*tracks)[idx][TagMode] = entry.Mode
		(*tracks)[idx][TagSpeechiness] = entry.Speechiness
		(*tracks)[idx][TagTempo] = entry.Tempo
		(*tracks)[idx][TagTimeSignature] = entry.TimeSignature
		(*tracks)[idx][TagValence] = entry.Valence
	}
}
