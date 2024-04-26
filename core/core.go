package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
)

func GetFavorites(client *spotify.Client, update chan<- string) ([]spotify.SavedTrack, error) {
	savedTrackPage, err := client.CurrentUsersTracks(context.Background(), spotify.Limit(50))
	total := savedTrackPage.Total
	savedTracks := make([]spotify.SavedTrack, 0, total)

	for err == nil {
		savedTracks = append(savedTracks, savedTrackPage.Tracks...)
		update <- fmt.Sprintf("favorites: %d%%", 100*len(savedTracks)/total)
		err = client.NextPage(context.Background(), savedTrackPage)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}
	return savedTracks, nil
}

func GetAudioFeatures(
	client *spotify.Client,
	favorites []spotify.SavedTrack,
	update chan<- string) ([]*spotify.AudioFeatures, error) {

	total := len(favorites)
	audioFeatures := make([]*spotify.AudioFeatures, 0, total)

	for i := 0; i <= total/100; i++ {
		ids := make([]spotify.ID, 0, 100)

		start := i * 100
		end := min(total, start+100)

		for _, entry := range favorites[start:end] {
			ids = append(ids, entry.ID)
		}

		current, err := client.GetAudioFeatures(context.Background(), ids...)
		if err != nil {
			return audioFeatures, err
		}

		audioFeatures = append(audioFeatures, current...)
		update <- fmt.Sprintf("audio features: %d%%", 100*len(audioFeatures)/total)
	}

	return audioFeatures, nil
}

func WriteToFile(
	tracks []map[string]interface{},
	fileName string) error {

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
