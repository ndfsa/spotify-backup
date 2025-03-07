package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os/user"

	"github.com/zalando/go-keyring"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var clientId = "INVALID"

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	port    = "8080"
	service = "spotify-backup"
)

func generateRandomString(n int) string {
	buffer := make([]byte, n)
	for i := range buffer {
		next, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err)
		}
		buffer[i] = charset[next.Int64()]
	}
	return string(buffer)
}

func GetClient() *spotify.Client {
	redirectURI := "http://localhost:" + port + "/callback"

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopePlaylistReadPrivate),
		spotifyauth.WithClientID(clientId))

	token := getToken(auth)

	return spotify.New(auth.Client(context.Background(), token))
}

func getToken(auth *spotifyauth.Authenticator) *oauth2.Token {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	secret, err := keyring.Get(service, user.Username)
	if err != nil {
		return getNewToken(auth, port, user.Username)
	}

	token, err := refreshToken(secret, auth, user.Username)
	if err != nil {
		return getNewToken(auth, port, user.Username)
	}

	return token
}

func refreshToken(
	secret string,
	auth *spotifyauth.Authenticator,
	user string,
) (*oauth2.Token, error) {
	token, err := auth.RefreshToken(context.Background(), &oauth2.Token{RefreshToken: secret})
	if err != nil {
		return nil, err
	}

	if err := keyring.Set(service, user, token.RefreshToken); err != nil {
		log.Fatal(err)
	}

	return token, nil
}

func getNewToken(auth *spotifyauth.Authenticator, port, user string) *oauth2.Token {
	state := generateRandomString(64)
	ch := make(chan *oauth2.Token)

	// PKCE variables
	codeVerifier := generateRandomString(128)
	sum := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			return
		}

		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
		}

		ch <- token

		fmt.Fprint(w, "Login successful!")
	})

	srv := &http.Server{Addr: ":" + port}
	go srv.ListenAndServe()

	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge))
	fmt.Println("Follow and login: ", url)

	token := <-ch
	srv.Shutdown(context.Background())

	err := keyring.Set(service, user, token.RefreshToken)
	if err != nil {
		log.Fatal(err)
	}

	return token
}
