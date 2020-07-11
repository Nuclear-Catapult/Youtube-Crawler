crawler: main.go parseHTML.go parseJSON.go cache
	go build main.go parseHTML.go parseJSON.go

cache:
	$(MAKE) -C DB-Cache
