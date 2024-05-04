package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

type DumpTrack struct {
	AddedAt  string   `json:"added_at"`
	Album    string   `json:"album"`
	Artists  []string `json:"artists"`
	Duration int      `json:"duration"`
	Explicit bool     `json:"explicit"`
	Id       string   `json:"id"`
	Name     string   `json:"name"`
}

type FullDumpTrack struct {
	Acousticness     float32 `json:"acousticness"`
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Instrumentalness float32 `json:"instrumentalness"`
	Key              int     `json:"key"`
	Liveness         float32 `json:"liveness"`
	Loudness         float32 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float32 `json:"speechiness"`
	Tempo            float32 `json:"tempo"`
	TimeSignature    int     `json:"time_signature"`
	Valence          float32 `json:"valence"`
	DumpTrack
}

func NewFromPlaylistTrack(track spotify.PlaylistTrack) DumpTrack {
	fullTrack := track.Track
	res := DumpTrack{
		AddedAt:  track.AddedAt,
		Album:    fullTrack.Album.Name,
		Artists:  make([]string, 0, len(fullTrack.Artists)),
		Duration: fullTrack.Duration,
		Explicit: fullTrack.Explicit,
		Id:       fullTrack.ID.String(),
		Name:     fullTrack.Name,
	}
	for _, artist := range fullTrack.Artists {
		res.Artists = append(res.Artists, artist.Name)
	}
	return res
}

func NewFromSavedTrack(track spotify.SavedTrack) DumpTrack {
	dump := DumpTrack{
		AddedAt:  track.AddedAt,
		Album:    track.Album.Name,
		Artists:  make([]string, 0, len(track.Artists)),
		Duration: track.Duration,
		Explicit: track.Explicit,
		Id:       track.ID.String(),
		Name:     track.Name,
	}
	for _, artist := range track.Artists {
		dump.Artists = append(dump.Artists, artist.Name)
	}
	return dump
}

func UpgradeDumpTrack(track DumpTrack, features *spotify.AudioFeatures) FullDumpTrack {
	return FullDumpTrack{
		DumpTrack:        track,
		Acousticness:     features.Acousticness,
		Danceability:     features.Danceability,
		Energy:           features.Energy,
		Instrumentalness: features.Instrumentalness,
		Key:              features.Key,
		Liveness:         features.Liveness,
		Loudness:         features.Loudness,
		Mode:             features.Mode,
		Speechiness:      features.Speechiness,
		Tempo:            features.Tempo,
		TimeSignature:    features.TimeSignature,
		Valence:          features.Valence,
	}
}

func GetFavorites(client *spotify.Client, update chan<- string) ([]DumpTrack, error) {
	savedTrackPage, err := client.CurrentUsersTracks(context.Background(), spotify.Limit(50))
	total := savedTrackPage.Total
	tracks := make([]DumpTrack, 0, total)

	for err == nil {
		for _, elem := range savedTrackPage.Tracks {
			tracks = append(tracks, NewFromSavedTrack(elem))
		}
		update <- fmt.Sprintf("favorites: %d%%", 100*len(tracks)/total)
		err = client.NextPage(context.Background(), savedTrackPage)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}
	return tracks, nil
}

func GetPlaylist(
	client *spotify.Client,
	id spotify.ID,
	update chan<- string,
) ([]DumpTrack, error) {
	playlist, err := client.GetPlaylist(context.Background(), id, spotify.Limit(50))
	if err != nil {
		return []DumpTrack{}, err
	}
	total := playlist.Tracks.Total
	tracks := make([]DumpTrack, 0, total)

	playlistItemPage := playlist.Tracks

	for err == nil {
		for _, elem := range playlistItemPage.Tracks {
			tracks = append(tracks, NewFromPlaylistTrack(elem))
		}
		update <- fmt.Sprintf("%s: %d%%", playlist.Name, 100*len(tracks)/total)
		err = client.NextPage(context.Background(), &playlistItemPage)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}
	return tracks, nil
}

func GetAudioFeatures(
	client *spotify.Client,
	tracks []DumpTrack,
	update chan<- string,
) ([]FullDumpTrack, error) {
	total := len(tracks)
	fullTracks := make([]FullDumpTrack, 0, total)

	for i := 0; i <= total/100; i++ {
		start := i * 100
		end := min(total, start+100)
		chunk := tracks[start:end]

		ids := make([]spotify.ID, 0, len(chunk))
		for _, elem := range chunk {
			ids = append(ids, spotify.ID(elem.Id))
		}

		features, err := client.GetAudioFeatures(context.Background(), ids...)
		if err != nil {
			return fullTracks, err
		}
		for idx, elem := range chunk {
			fullTracks = append(fullTracks, UpgradeDumpTrack(elem, features[idx]))
		}

		update <- fmt.Sprintf("audio features: %d%%", 100*len(fullTracks)/total)
	}

	return fullTracks, nil
}

func WriteToFile(
	tracks any,
	fileName string,
) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	jsonEnc := json.NewEncoder(w)
	if err := jsonEnc.Encode(&tracks); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
