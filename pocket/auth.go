package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Auth(cKey string) (token string, err error) {
	rToken, err := requestToken(cKey)
	if err != nil {
		err = fmt.Errorf("couldn't get request token: %v", err)
	}
	fmt.Println(rToken)
	return
}

func requestToken(cKey string) (rToken string, err error) {
	data, err := json.Marshal(map[string]string{
		"consumer_key": cKey,
		"redirect_url": "http://sellweek.github.io/pockettxt",
	})
	if err != nil {
		err = fmt.Errorf("couldn't marshal data: %v", err)
		return
	}

	r, err := http.NewRequest("POST", "http://getpocket.com/v3/oauth/request", bytes.NewBuffer(data))
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't read response: %v", err)
		return
	}

	m := make(map[string]string, 1)
	err = json.Unmarshal(body, &m)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal response: %v", err)
		return
	}

	rToken = m["code"]
	return
}
