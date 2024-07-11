package spycat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	catsApiUrl     string = "https://api.thecatapi.com/v1/breeds"
	breedFieldName string = "name"
)

type Breed string

type Validator interface {
	Validate(Breed) error
}

type catValidator struct {
	// no mutex needed because this is supposed to be read only
	validBreeds map[Breed]struct{}
}

func NewCatValidator() (Validator, error) {
	var resp []map[string]interface{}

	httpResp, err := http.Get(catsApiUrl)
	if err != nil {
		return nil, fmt.Errorf("http request error: %s", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read http response body: %s", err)
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal catapi response: %s", err)
	}

	cv := catValidator{validBreeds: make(map[Breed]struct{})}
	for _, v := range resp {
		breedName, ok := v[breedFieldName].(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert %s to string", breedFieldName)
		}
		cv.validBreeds[Breed(breedName)] = struct{}{}
	}

	return &cv, nil
}

// returns list of possible breeds
func (cv *catValidator) Validate(breed Breed) error {
	_, exists := cv.validBreeds[breed]
	if !exists {
		return fmt.Errorf("breed %s is not recognized", breed)
	}
	return nil
}
