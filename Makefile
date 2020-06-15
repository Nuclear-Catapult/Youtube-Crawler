crawler: crawler.go parseHTML.go cache
	go build crawler.go parseHTML.go

cache:
	$(MAKE) -C ID-Cache
