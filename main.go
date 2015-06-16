package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	//	WIKI_PREFIX       = "http://en.wikipedia.org"
	WIKI_PREFIX       = "http://10.102.44.202"
	BEGIN_OFFSET      = 10000
	LINKS_AT_ONCE     = 50
	NUM_OF_ITERATIONS = 100

	RESULT_FILE = "links.txt"
)

func GrabLinks(wiki string, clinks chan []string) {
	doc, err := goquery.NewDocument(wiki)
	if err != nil {
		log.Fatal(err)
	}

	links := make([]string, 0, LINKS_AT_ONCE)
	oles := doc.Find("ol")

	oles.Find("li").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Eq(1).Attr("href")
		//		fmt.Printf("%s\n", link)
		links = append(links, link)
	})
	clinks <- links
}

func ScrapeAllWikis() {
	clinks := make(chan []string, NUM_OF_ITERATIONS)
	for idx := 0; idx < NUM_OF_ITERATIONS; idx++ {
		offset := LINKS_AT_ONCE*idx + BEGIN_OFFSET
		//"https://en.wikipedia.org/w/index.php?title=Special:LongPages&limit=5000&offset=0"
		//		link := fmt.Sprintf("%s/index.php?title=Special:LongPages&limit=%d&offset=%d",
		link := fmt.Sprintf("%s/index.php?title=Special:ShortPages&limit=%d&offset=%d",
			WIKI_PREFIX, LINKS_AT_ONCE, offset)
		//		fmt.Printf("Grabbing: %s\n", link)
		go GrabLinks(link, clinks)
	}

	f, err := os.OpenFile(RESULT_FILE, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for idx := 0; idx < NUM_OF_ITERATIONS; idx++ {
		select {
		case links := <-clinks:
			fmt.Printf("-")
			WriteLinksToFile(links, f)
		}
	}
	fmt.Printf("\n")
	fmt.Printf("fin\n")
}

func WriteLinksToFile(links []string, file *os.File) {
	for _, link := range links {
		//		fmt.Printf("%s\n", link)
		file.WriteString(WIKI_PREFIX + link + "\n")
	}
}

func main() {
	ScrapeAllWikis()
}
