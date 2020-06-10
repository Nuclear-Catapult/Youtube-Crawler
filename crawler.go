package main

import (
	"fmt"
	"log"
	cache "github.com/Nuclear-Catapult/NCcrawl-YT/ID-Cache"
    "github.com/PuerkitoBio/goquery"
)

func main() {
	seed_ID := "hsWr_JWTZss"
	cache.Insert(seed_ID)

	crawler()
}

func crawler() {
	for id := cache.Next(); id != ""; id = cache.Next() {
		doc, err := goquery.NewDocument("https://www.youtube.com/watch?v=" + id)
		checkErr(err)
		doc.Find(".content-link.spf-link").Each(func(index int, item *goquery.Selection){
			link, err := item.Attr("href")
			if err == false {
				fmt.Println("Error: No href attribute found")
			}
			cache.Insert(string(link[len(link)-11:len(link)]))
		})
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Uh Oh..")
		log.Fatal(err)
	}
}
