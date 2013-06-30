package pocket

import (
	"encoding/json"
	"fmt"
	"github.com/toqueteos/webbrowser"
	"net/http"
)

const (
	redirURI = "http://127.0.0.1:50212"
)

func Auth(cKey string) (token string, err error) {
	fmt.Println("Beginning authentication")
	rToken, err := requestToken(cKey)
	if err != nil {
		return
	}

	url := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", rToken, redirURI)
	if err = webbrowser.Open(url); err == nil {
		c := make(chan bool)
		fmt.Println("Communicating with Pocket.")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Got a response from Pocket.")
			fmt.Fprint(w, "You can close this window now.")
			c <- true
		})
		go func() {
			if err := http.ListenAndServe(":50212", nil); err != nil {
				panic(err)
			}
		}()

		fmt.Println("Waiting for response.")
		<-c
	} else {
		fmt.Println("Please open the following URL in your browser, log in and then press enter")
		fmt.Println(url)
		var str string
		fmt.Scanf("%s", &str)
	}

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
