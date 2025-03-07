include .env

VERSION := $(shell date +%y.%m.%d).$(shell git rev-parse --short HEAD)
LDFLAGS := -s -w -X spotify-backup/auth.clientId=${SPOTIFY_ID} -X main.versionId=${VERSION}

all: build

build:
	@go build -ldflags="${LDFLAGS}"

clean:
	@rm spotify-backup
