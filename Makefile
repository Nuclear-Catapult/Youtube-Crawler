crawler: main.go parseHTML.go cache
	go build main.go parseHTML.go

cache:
	$(MAKE) -C DB-Cache
