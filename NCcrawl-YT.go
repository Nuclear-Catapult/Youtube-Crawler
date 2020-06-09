package main

import cache "github.com/Nuclear-Catapult/NCcrawl-YT/ID-Cache"
import "fmt"

func main() {
	cache.Insert("GGsf8444444")
	cache.Insert("lksf8440444")
	cache.Insert("lksf8440444")
	cache.Insert("siwi2444444")
	fmt.Println(cache.Next())
	fmt.Println(cache.Next())
	fmt.Println(cache.Next())
	fmt.Println(cache.Next())
}
