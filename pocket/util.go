package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func requestChain(data map[string]string, url string) (response []byte, err error) {
	jsData, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("couldn't marshal data: %v", err)
		return
	}

	r, err := http.NewRequest("POST", "http://getpocket.com"+url, bytes.NewBuffer(jsData))
	if err != nil {
		err = fmt.Errorf("couldn't create request: %v", err)
		return
	}

	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("X-Accept", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		err = fmt.Errorf("error when making request: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Pocket returned error: %s", resp.Header.Get("X-Error"))
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't read response: %v", err)
		return
	}

	return
}
