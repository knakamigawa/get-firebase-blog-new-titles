package main

import (
	"bytes"
	"fmt"
	"get-firebase-blog-new-titles/client"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "https://firebase.blog/"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(body)

	node, err := html.Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	var collection []string
	findH3(node, &collection)

	aiClient := client.ProvideAPIClient()

	if len(collection) > 0 {
		for _, line := range collection {
			fmt.Println(line)
			res, err := aiClient.Request(line)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res)
		}
	}
}

func findH3(node *html.Node, collection *[]string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.DataAtom == atom.H3 {
				*collection = append(*collection, c.FirstChild.FirstChild.Data)
			}
			findH3(c, collection)
		}
	}
}
