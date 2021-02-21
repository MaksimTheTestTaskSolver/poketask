package http

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"

	"github.com/disintegration/imaging"
)

func Get(url string, respDestination interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't make GET request to %s: %w", url, err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read the response body from %s: %w", url, err)
	}

	if resp.StatusCode > 400 {
		return fmt.Errorf("error code in response from %s: %d - %s\n", url, resp.StatusCode, string(respBody))
	}

	err = json.Unmarshal(respBody, respDestination)
	if err != nil {
		return fmt.Errorf("can't unmarshal body from %s: %w", url, err)
	}

	return nil
}

func GetImage(url string) (image image.Image, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't make GET request to %s: %w", url, err)
	}

	if resp.StatusCode > 400 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't read the response body from %s: %w", url, err)
		}
		return nil, fmt.Errorf("error code in response from url %s: %d - %s\n", url, resp.StatusCode, string(respBody))
	}

	image, err = imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't decode image from response body from %s: %w", url, err)
	}

	return image, nil
}

