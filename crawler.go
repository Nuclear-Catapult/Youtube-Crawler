package main

import (
	"fmt"
	"log"
	cache "github.com/Nuclear-Catapult/Youtube-Crawler/ID-Cache"
	b64 "github.com/Nuclear-Catapult/Youtube-Crawler/ytbase64"
    "github.com/PuerkitoBio/goquery"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	seed_ID := "hsWr_JWTZss"
	cache.Insert(b64.Decode64(seed_ID))
	c := make(chan []interface{})

	go inserter(c)
	crawler(c)
}

func inserter(c chan []interface{}) {
	db, err := sql.Open("sqlite3", "./yt-videos.db")
	checkErr(err)

	table_stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS video (
	video_id INTEGER(64) PRIMARY KEY,
	rec_1 INTEGER(64) NOT NULL,
	rec_2 INTEGER(64) NOT NULL,
	rec_3 INTEGER(64) NOT NULL,
	rec_4 INTEGER(64) NOT NULL,
	rec_5 INTEGER(64) NOT NULL,
	rec_6 INTEGER(64) NOT NULL,
	rec_7 INTEGER(64) NOT NULL,
	rec_8 INTEGER(64) NOT NULL,
	rec_9 INTEGER(64) NOT NULL,
	rec_10 INTEGER(64) NOT NULL,
	rec_11 INTEGER(64) NOT NULL,
	rec_12 INTEGER(64) NOT NULL,
	rec_13 INTEGER(64) NOT NULL,
	rec_14 INTEGER(64) NOT NULL,
	rec_15 INTEGER(64) NOT NULL,
	rec_16 INTEGER(64) NOT NULL,
	rec_17 INTEGER(64) NOT NULL,
	rec_18 INTEGER(64) NOT NULL);`)
	checkErr(err)
	table_stmt.Exec()

	insert_stmt, err := db.Prepare(`INSERT INTO video
	(video_id, rec_1, rec_2, rec_3, rec_4, rec_5, rec_6, rec_7, rec_8, rec_9,
	rec_10, rec_11, rec_12, rec_13, rec_14, rec_15, rec_16, rec_17, rec_18)
	values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	checkErr(err)

	for {
		_, err := insert_stmt.Exec(<-c...)
		checkErr(err)
	}
}

func crawler(c chan []interface{}) {
	for id := cache.Next(); id != 0; id = cache.Next() {
		var	rec_count int
		row := []interface{}{int64(id)} // This will be inserted into yt-videos.db
		doc, err := goquery.NewDocument("https://www.youtube.com/watch?v=" + b64.Encode64(id))
		checkErr(err)
		rec_sel := doc.Find(".content-link.spf-link")
		if rec_sel.Length() < 18 {
		// For some reason, a valid YT webage varies with its initial recommendation count. Downloading a webpage
		// may yield 22 recommendations, and downloading the same page again usually results in a
		// different count. A minority of pages have less than 18, of these we'll insert back into the queue to
		// try again later.
			 cache.TryAgainLater(id)
			 continue
		}
		rec_sel.EachWithBreak(func(index int, item *goquery.Selection) bool {
			link, err := item.Attr("href")
			if err == false {
				fmt.Println("Error: No href attribute found")
			}
			rec_id := b64.Decode64(string(link[len(link)-11:len(link)]))
			cache.Insert(rec_id)
			row = append(row, int64(rec_id))
			rec_count++
			if rec_count == 18 {
				return false
			}
			return true
		})
		c <- row
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Uh Oh..")
		log.Fatal(err)
	}
}
