package main

import (
	"fmt"
	cache "github.com/Nuclear-Catapult/Youtube-Crawler/DB-Cache"
	b64 "github.com/Nuclear-Catapult/Youtube-Crawler/ytbase64"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
)

func extractNumber(str string, doc *goquery.Document) int64 {
	if str == "" {
		return -1
	}
	index := strings.IndexByte(str, ' ')
	if index != -1 {
		str = str[:index]
	}
	if str == "No" {
		// This should mean the video has "No views"
		return 0
	}
	views, err := strconv.ParseInt(strings.ReplaceAll(str, ",", ""), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return views
}

func ParseHTML(doc *goquery.Document, id int64, row []interface{}, video_chan chan []interface{}) {
	views_slice := doc.Find(".watch-view-count").Text()
	row = append(row, extractNumber(views_slice, doc))

	likes_slice := doc.Find(".like-button-renderer-like-button-unclicked > span:nth-child(1)").Text()
	dislikes_slice := doc.Find(".like-button-renderer-dislike-button-unclicked > span:nth-child(1)").Text()
	row = append(row, extractNumber(likes_slice, doc))
	row = append(row, extractNumber(dislikes_slice, doc))

	channel_id, status := doc.Find(`[itemprop="channelId"]`).Attr("content")
	if status == false {
		f, _ := os.Create("failed_fluke.html")
		html, _ := doc.Html()
		f.WriteString(html)
		log.Fatal("content empty")
	}
	lhalf := b64.Decode64(channel_id[2:])
	rhalf := b64.Decode64(channel_id[13:])

	cache.InsertChannel(lhalf, rhalf)

	rec_sel := doc.Find(".content-link.spf-link")
	if rec_sel.Length() < 18 {
		// For some reason, a valid YT webage varies with its initial recommendation count. Downloading a webpage
		// may yield 22 recommendations, and downloading the same page again usually results in a
		// different count. A minority of pages have less than 18, of these we'll insert back into the queue to
		// try again later.
		cache.TryAgainLater(id)
		return
	}

	var rec_count int
	rec_sel.EachWithBreak(func(index int, item *goquery.Selection) bool {
		link, err := item.Attr("href")
		if err == false {
			fmt.Println("Error: No href attribute found")
		}
		rec_id := b64.Decode64(string(link[len(link)-11 : len(link)]))
		cache.Insert(rec_id)
		row = append(row, rec_id)
		rec_count++
		if rec_count == 18 {
			return false
		}
		return true
	})
	video_chan <- row
}
