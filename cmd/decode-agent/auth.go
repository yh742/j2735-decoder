package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type basicAuth struct {
	username string
	password string
}

func checkBasicHTTPAuth(r *http.Request, auth basicAuth) bool {
	if auth.username != "" && auth.password != "" {
		u, p, ok := r.BasicAuth()
		if !ok {
			log.Error().Msg("auth information not supplied")
			return false
		}
		if strings.TrimSpace(u) == auth.username && strings.TrimSpace(p) == auth.password {
			return true
		}
		log.Error().Msg("auth information is wrong")
		return false
	}
	log.Debug().Msg("HTTP Auth not performed, missing username and/or password")
	return true
}

func parseAuthFiles(path string) basicAuth {
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal().Err(err).Msg("error occured accessing http secrets")
		}
		bytes, err := ioutil.ReadAll(file)
		sArr := strings.Split(string(bytes), "\n")
		if len(sArr) < 2 {
			log.Fatal().Err(err).Msg("password file format is incorrect")
		}
		user := strings.TrimSpace(sArr[0])
		psw := strings.TrimSpace(sArr[1])
		if user == "" || psw == "" {
			log.Fatal().Err(err).Msg("username or password cannot be empty")
		}
		return basicAuth{
			username: user,
			password: psw,
		}
	}
	// no auth required
	return basicAuth{}
}
