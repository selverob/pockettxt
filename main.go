package main

import (
	"flag"
	"fmt"
	"github.com/sellweek/pockettxt/article"
	"github.com/sellweek/pockettxt/diffbot"
	"github.com/sellweek/pockettxt/pocket"
	"io"
	"os"
)

const (
	cKey   = "14774-3dbc1ccffc2398d5cf6cefe1"
	dToken = "384e44301a3ffe07396c81781c46b7e9"
)

var (
	fToken = flag.String("aToken", "", "Authorization token provided by Pocket, if you have one.")
)

func main() {
	flag.Parse()

	var (
		aToken string
		err    error
	)

	if *fToken == "" {
		aToken, err = pocket.Auth(cKey)
		if err != nil {
			panic(err)
		}
	} else {
		aToken = *fToken
	}

	fmt.Println("Got access token: ", aToken)

	urls, err := pocket.URLs(cKey, aToken)
	if err != nil {
		panic(err)
	}

	as, err := diffbot.Articles(urls, dToken)
	if err != nil {
		panic(err)
	}

	err = writeArticles("pocket.txt", as)
}

func writeArticles(filename string, as []article.Article) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		err = fmt.Errorf("couldn't create output file: %v", err)
		return
	}

	_, err = io.Copy(f, article.PrintArticles(as))
	if err != nil {
		err = fmt.Errorf("couldn't write articles to a file: %v", err)
		return
	}

	return
}