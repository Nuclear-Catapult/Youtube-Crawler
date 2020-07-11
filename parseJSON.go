package main

import (
	cache "github.com/Nuclear-Catapult/Youtube-Crawler/DB-Cache"
	b64 "github.com/Nuclear-Catapult/Youtube-Crawler/ytbase64"
	"encoding/json"
	"bytes"
	"log"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type struct3 struct {
	Results struct {
		Results struct {
			Contents []map[string]map[string]interface{}
		}
	}
	SecondaryResults map[string]interface{}
}

type struct2 struct {
	 TwoColumnWatchNextResults struct3
}

type root struct {
	Contents struct2
}

func ParseJSON(doc *goquery.Document, video_row []interface{}, id int64) bool {
	jsn := doc.Find("body > script").Eq(-4).Text()

	jsn = jsn[31:]
	offset := strings.IndexByte(jsn, '\n')

	jsn = jsn[:offset-1]

	j := []byte(jsn)
	status := json.Valid(j)

	if status == false {
	//	printJSON(j)
		log.Fatal("Invalid JSON format")
	}

	var r root
	err := json.Unmarshal(j, &r)
	if err != nil {
		log.Fatal(err)
	}

	contents := r.Contents.TwoColumnWatchNextResults.Results.Results.Contents
	if contents == nil {
		// Most likely a private video
		return true
	}
	contents_index := 0
	primary := contents[0]["videoPrimaryInfoRenderer"]
	if primary == nil {
		primary = contents[1]["videoPrimaryInfoRenderer"]
		contents_index++
	}

	{// get title
		title := primary["title"].
		(map[string]interface{})["runs"].
		([]interface{})[0].
		(map[string]interface{})["text"]

		video_row = append(video_row, title)
	}

	{// get views
		views := primary["viewCount"]
		// If views == nil, then video is either a movie or livestream. We don't want this in our database
		if views == nil {
			return true
		}
		views = views.(map[string]interface{})["videoViewCountRenderer"].
		(map[string]interface{})["viewCount"].
		(map[string]interface{})["simpleText"]
		if views == nil {
			return true
		}

		video_row = append(video_row, ExtractNumber(views.(string)))
	}

	{// get ratings
		ratings := primary["videoActions"].
		(map[string]interface{})["menuRenderer"].
		(map[string]interface{})["topLevelButtons"]

		likes := ratings.([]interface{})[0].
		(map[string]interface{})["toggleButtonRenderer"].
		(map[string]interface{})["defaultText"].
		(map[string]interface{})["accessibility"]

		if likes == nil {
			// ratings are disabled
			video_row = append(video_row, -1) // likes = -1
			video_row = append(video_row, -1) // dislikes = -1
		} else {
			likes = likes.(map[string]interface{})["accessibilityData"].
			(map[string]interface{})["label"]

			video_row = append(video_row, ExtractNumber(likes.(string)))

			dislikes := ratings.([]interface{})[1].
			(map[string]interface{})["toggleButtonRenderer"].
			(map[string]interface{})["defaultText"].
			(map[string]interface{})["accessibility"].
			(map[string]interface{})["accessibilityData"].
			(map[string]interface{})["label"]

			video_row = append(video_row, ExtractNumber(dislikes.(string)))
		}
	}

	/*
	{// get channel name, ID, and subscriber count
		subscriber_count := contents[1]["videoSecondaryInfoRenderer"]["owner"].
		(map[string]interface{})["videoOwnerRenderer"].
		(map[string]interface{})["subscriberCountText"].
		(map[string]interface{})["runs"].
		([]interface{})[0].
		(map[string]interface{})["text"]

		channel := contents[1]["videoSecondaryInfoRenderer"]["owner"].
		(map[string]interface{})["videoOwnerRenderer"].
		(map[string]interface{})["title"].
		(map[string]interface{})["runs"].
		([]interface{})[0]

		channel_id := channel.(map[string]interface{})["navigationEndpoint"].
		(map[string]interface{})["browseEndpoint"].
		(map[string]interface{})["browseId"]

		lhalf := b64.Decode64(channel_id.(string)[2:])
		rhalf := b64.Decode64(channel_id.(string)[13:])

		if subscriber_count != "" && cache.InsertChannel(lhalf, rhalf) == true {
			channel_name := channel.(map[string]interface{})["text"]

			if channel_name == "YouTube Movies" {
				log.Fatal("This should be impossible. A movie made it down to channel handling")
			}

			channel_row := []interface{}{lhalf, rhalf}
			channel_row = append(channel_row, html.EscapeString(channel_name.(string)))
			channel_row  = append(channel_row, GetSubs(subscriber_count.(string)))
			c <- channel_row
			fmt.Printf("%v %v %v %v\n", lhalf, rhalf, channel_name.(string), GetSubs(subscriber_count.(string)))
		}

		video_row = append(video_row, lhalf)
		video_row = append(video_row, rhalf)
	}
	*/

	rec_check :=  r.Contents.TwoColumnWatchNextResults.SecondaryResults["secondaryResults"]
	if rec_check == nil || rec_check.(map[string]interface{})["results"] == nil {
		zero_recommendations(&video_row)
		return true
	}

	recommendations := r.Contents.TwoColumnWatchNextResults.SecondaryResults["secondaryResults"].
	(map[string]interface{})["results"].([]interface{})

	// The first recommendation is more deep since it's in autoplay
	rec_id := recommendations[0].(map[string]interface{})["compactAutoplayRenderer"]

	if rec_id == nil {
		zero_recommendations(&video_row)
		return true
	}

	rec_id = rec_id.(map[string]interface{})["contents"].
	([]interface{})[0].
	(map[string]interface{})["compactVideoRenderer"].
	(map[string]interface{})["videoId"]

	rec_id_64 := b64.Decode64(rec_id.(string))
	cache.Insert(rec_id_64)
	video_row = append(video_row, rec_id_64)

	content_count := len(recommendations)
	max := 18
	for i := 1; i < max; i++ {
		if i == content_count {
			return false
		}
		rec_id = recommendations[i]

		rec_id = rec_id.(map[string]interface{})["compactVideoRenderer"]

		if rec_id == nil {
			max++
			continue
		}

		rec_id = rec_id.(map[string]interface{})["videoId"]
		rec_id_64 = b64.Decode64(rec_id.(string))
		cache.Insert(rec_id_64)
		video_row = append(video_row, rec_id_64)
	}
	v <- video_row
	return true
}

func printJSON(jsn []byte) {
	f, _ := os.Create("dump.json")
	defer f.Close()
	var out bytes.Buffer
	err := json.Indent(&out, jsn, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	out.WriteTo(f)
}

func zero_recommendations(video_row *[]interface{}) {
	for i := 0; i < 18; i++ {
		*video_row = append(*video_row, 0)
	}
	v <- *video_row
}
