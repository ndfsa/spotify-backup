package core

import (
	"bufio"
	"context"
	"os"

	"github.com/zmb3/spotify/v2"
)

func GetFavorites(client *spotify.Client) ([]spotify.SavedTrack, error) {
	savedTracks := make([]spotify.SavedTrack, 0, 50)

	savedTrackPage, err := client.CurrentUsersTracks(context.Background(), spotify.Limit(50))

	for err == nil {
		savedTracks = append(savedTracks, savedTrackPage.Tracks...)
		err = client.NextPage(context.Background(), savedTrackPage)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}
	return savedTracks, nil
}

func WriteToFile(
	savedTracks []spotify.SavedTrack,
	fields uint,
	encoder SavedTracksEncoder,
	fileName string) error {

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	if err := encoder.Encode(savedTracks, fields, w); err != nil {
		return err
	}

	return nil
}
