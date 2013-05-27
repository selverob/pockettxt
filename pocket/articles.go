package pocket

import (
	"encoding/json"
	"fmt"
)

func URLs(cKey, aToken string) (urls []string, err error) {
	fmt.Println("Getting article list from Pocket")

	data := map[string]string{
		"consumer_key": cKey,
		"access_token": aToken,
		"contentType":  "article",
		"detailType":   "simple",
	}
	resp, err := requestChain(data, "/v3/get")
	if err != nil {
		err = fmt.Errorf("couldn't get articles from Pocket: %v", err)
		return
	}

	as := make(map[string]interface{})
	err = json.Unmarshal(resp, &as)
	if err != nil {
		err = fmt.Errorf("couldn't unmarshal JSON: %v", err)
		return
	}

	articleList := as["list"].(map[string]interface{})
	urls = make([]string, 0, len(articleList))

	for _, a := range articleList {
		urls = append(urls, a.(map[string]interface{})["resolved_url"].(string))
	}

	return
}
