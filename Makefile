crawler: crawler.go cache
	go build crawler.go

cache:
	$(MAKE) -C ID-Cache
