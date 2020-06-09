package main

import cache "github.com/Nuclear-Catapult/NCcrawl-YT/ID-Cache"
import "fmt"

func main() {
	seed_ID := "hsWr_JWTZss"
	cache.Insert(seed_ID)

	crawler()
}

func crawler() {
	for id := cache.Next(); id != ""; id = cache.Next() {
		fmt.Println(id)
	}
}
