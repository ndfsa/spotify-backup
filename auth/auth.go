package auth

import (
	"fmt"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/joho/godotenv"
)

const (
	PORT        = "8080"
	redirectURI = "http://localhost:" + PORT + "/callback"
)

var (
	_    = godotenv.Load()
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopePlaylistReadPrivate),
	)
	state = "backup_state"
)

func SetupAuth(ch chan<- *spotify.Client) {
	// setup an http server to receive OAuth token
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// get token from the request
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
		}

		// create new spotify client
		client := spotify.New(auth.Client(r.Context(), token))

		fmt.Fprint(w, "Login successful!")

		// send client to chanel
		ch <- client
	})

	// launch callback server in background task
	go func() {
		http.ListenAndServe(":"+PORT, nil)
	}()

	url := auth.AuthURL(state)
	fmt.Println("Follow and login: ", url)
}
