package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	strip "github.com/grokify/html-strip-tags-go"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

func main() {
	response, err := http.Get("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	rss := new(Rss)
	buffer := bytes.NewBuffer(XMLdata)
	decoded := xml.NewDecoder(buffer)
	err = decoded.Decode(rss)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Title : %s\n", rss.Channel.Title)
	fmt.Printf("Description : %s\n", rss.Channel.Description)
	fmt.Printf("Link : %s\n", rss.Channel.Link)
	for i, elem := range rss.Channel.Items {

		fmt.Printf("[%d] item title : %s\n", i, elem.Title)
		fmt.Printf("[%d] item description : %s\n", i, strip.StripTags(elem.Description))
		fmt.Printf("[%d] item link : %s\n\n", i, elem.Link)
		fmt.Printf("[%d] item pubDate : %s\n\n", i, elem.PubDate)
	}
	/*
		fmt.Printf("Title : %s\n", rss.Channel.Title)
		fmt.Printf("Description : %s\n", rss.Channel.Description)
		fmt.Printf("Link : %s\n", rss.Channel.Link)

		total := len(rss.Channel.Items)

		fmt.Printf("Total items : %v\n", total)

		for i := 0; i < total; i++ {
			fmt.Printf("[%d] item title : %s\n", i, rss.Channel.Items[i].Title)
			fmt.Printf("[%d] item description : %s\n", i, rss.Channel.Items[i].Description)
			fmt.Printf("[%d] item link : %s\n\n", i, rss.Channel.Items[i].Link)
		}
	*/
}
