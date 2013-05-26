package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	redirURI = "http://sellweek.github.io/pockettxt"
)

func Auth(cKey string) (token string, err error) {
	fmt.Println("Beginning authentication")
	rToken, err := requestToken(cKey)
	if err != nil {
		return
	}

	url := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", rToken, redirURI)
	fmt.Println("Please open the following URL in your browser, log in and then press enter")
	fmt.Println(url)

	var str string
	fmt.Scanf("%s", &str)
	fmt.Println("Continuing authetication")

	token, err = accessToken(cKey, rToken)
	if err != nil {
		return
	}
	return
}

func requestToken(cKey string) (rToken string, err error) {
	data := map[string]string{
		"consumer_key": cKey,
		"redirect_uri": redirURI,
	}

	r, err := requestChain(data, "/v3/oauth/request")
	if err != nil {
		fmt.Errorf("couldn't get a request token: %v", err)
	}

	rToken = r["code"]
	return
}

func accessToken(cKey, rToken string) (aToken string, err error) {
	data := map[string]string{
		"consumer_key": cKey,
		"code":         rToken,
	}
	r, err := requestChain(data, "/v3/oauth/authorize")
	if err != nil {
		err = fmt.Errorf("couldn't get the access token: %v", err)
	}

	aToken = r["access_token"]
	return
}

func requestChain(data map[string]string, url string) (response map[string]string, err error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't read response: %v", err)
		return
	}

	response = make(map[string]string)
	err = json.Unmarshal(body, &response)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal response: %v", err)
		return
	}

	return
}
