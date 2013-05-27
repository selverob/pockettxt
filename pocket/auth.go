package pocket

import (
	"encoding/json"
	"fmt"
	"github.com/toqueteos/webbrowser"
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
	if err = webbrowser.Open(url); err == nil {
		fmt.Println("Login page should have opened in your browser")
		fmt.Println("Please log in and then come back and press enter")
	} else {
		fmt.Println("Please open the following URL in your browser, log in and then press enter")
		fmt.Println(url)
	}

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

	body, err := requestChain(data, "/v3/oauth/request")
	if err != nil {
		fmt.Errorf("couldn't get a request token: %v", err)
	}

	r := make(map[string]string)
	err = json.Unmarshal(body, &r)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal response: %v", err)
		return
	}

	rToken = r["code"]
	return
}

func accessToken(cKey, rToken string) (aToken string, err error) {
	data := map[string]string{
		"consumer_key": cKey,
		"code":         rToken,
	}
	body, err := requestChain(data, "/v3/oauth/authorize")
	if err != nil {
		err = fmt.Errorf("couldn't get the access token: %v", err)
	}

	r := make(map[string]string)
	err = json.Unmarshal(body, &r)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal response: %v", err)
		return
	}

	aToken = r["access_token"]
	return
}
