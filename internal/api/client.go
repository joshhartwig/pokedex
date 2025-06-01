package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joshhartwig/pokedex/pkg/models"
)

// fetchFromCache looks for url string in cache and returns associated data
// if the url is not found in the cache it will download the data and add the
// url and associated data encoded in json format to the map and return it
func FetchFromCache(c *models.Config, url string, v any) error {
	// try to find the url in cache 1st
	found, ok := c.Cache.Entries[url]
	if !ok { // if we did not find it, add a new cache entry with the data
		client := &http.Client{}
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// convert the resp.body to byte slide
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// add to cache
		c.Cache.Add(url, data)

		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&v); err != nil {
			fmt.Println("error decoding bytes", err)
			return err
		}
		return nil
	}

	// we found the url in the cache, return the data
	if err := json.NewDecoder(bytes.NewReader(found.Val)).Decode(&v); err != nil {
		fmt.Println("error decoding:", err)
		return err
	}

	return nil
}
