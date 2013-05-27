package diffbot

import (
	"encoding/json"
	"fmt"
	"github.com/sellweek/pockettxt/article"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Articles(urls []string, token string) (as []article.Article, err error) {
	fmt.Println("Requesting articles from Diffbot")

	reqData := make([]map[string]string, len(urls))
	for i, artUrl := range urls {
		m := make(map[string]string)
		m["method"] = "GET"
		m["relative_url"] = makeReqUrl(artUrl, token)
		reqData[i] = m
	}

	data, err := json.Marshal(reqData)
	if err != nil {
		err = fmt.Errorf("couldn't marshal JSON: %v", err)
		return
	}

	resp, err := makeRequest(data, token)
	if err != nil {
		return
	}

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't read response: %v", err)
		return
	}

	fmt.Println("Parsing response from Diffbot")
	as, err = decodeResponse(respJSON)

	if err != nil {
		err = fmt.Errorf("couldn't decode response: %v", err)
		return
	}

	return
}

func makeReqUrl(artUrl, token string) string {
	return fmt.Sprintf("/api/article?token=%s&url=%s",
		url.QueryEscape(token), url.QueryEscape(artUrl))
}

func makeRequest(batch []byte, token string) (r *http.Response, err error) {
	data := url.Values{}
	data.Add("token", token)
	data.Add("batch", string(batch))

	req, err := http.NewRequest("POST", "http://www.diffbot.com/api/batch?"+data.Encode(), nil)
	if err != nil {
		err = fmt.Errorf("couldn't form a request: %v", err)
		return
	}

	r, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("couldn't make a request to Diffbot: %v", err)
		return
	}

	if r.StatusCode != 200 {
		err = fmt.Errorf("server didn't respond with 200 but a %d", r.StatusCode)
	}

	return
}

func decodeResponse(js []byte) (as []article.Article, err error) {
	raw := make([]map[string]interface{}, 0)
	err = json.Unmarshal(js, &raw)
	if err != nil {
		err = fmt.Errorf("could not decode JSON: %v", err)
		return
	}

	as = make([]article.Article, 0)

	for _, r := range raw {
		rawArt := r["body"]
		artData := make(map[string]interface{})
		err = json.Unmarshal([]byte(rawArt.(string)), &artData)
		if err != nil {
			err = fmt.Errorf("couldn't decode article: %v", err)
			return
		}

		title, _ := artData["title"].(string)
		author, _ := artData["author"].(string)
		text, _ := artData["text"].(string)
		date, _ := artData["date"].(string)
		URL, _ := artData["url"].(string)

		a := article.Article{
			Title:  title,
			Author: author,
			Text:   text,
			Date:   date,
			URL:    URL,
		}
		as = append(as, a)
	}
	return

}
