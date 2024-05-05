package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var clientId = "INVALID"

func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buffer := make([]byte, 0, n)
	for ; n > 0; n-- {
		next, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		buffer = append(buffer, charset[next.Int64()])
	}
	return string(buffer)
}

func GetClient() *spotify.Client {
	// setup callback url
	port := "8080"
	redirectURI := "http://localhost:" + port + "/callback"

	// setup auth client with id
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopePlaylistReadPrivate),
		spotifyauth.WithClientID(clientId))

	// set state for OAuth workflow
	state := GenerateRandomString(64)

	// setup channels for execution
	ch := make(chan *spotify.Client)

	// setup PKCE variables
	codeVerifier := GenerateRandomString(128)
	sum := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])

	// setup an http server to receive OAuth token
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// get token from the request
		token, err := auth.Token(r.Context(), state, r,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			return
		}

		// verify the state
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
		}

		// create new spotify client
		client := spotify.New(auth.Client(r.Context(), token))

		fmt.Fprint(w, "Login successful!")

		// send client to chanel
		ch <- client
	})

	srv := &http.Server{Addr: ":" + port}

	// launch callback server in background task
	go srv.ListenAndServe()

	// show the authorization url to the user
	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge))
	fmt.Println("Follow and login: ", url)

	// wait for a client
	client := <-ch

	// shutdown the server
	srv.Shutdown(context.Background())

	return client
}
