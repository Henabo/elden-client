package gxios

import (
	"io"
	"log"
	"net/http"
)

func GET(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	return respBytes
}
