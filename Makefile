test:
	# go test -v -bench=. -benchmem -run=. -cpu=1,2 ./...
	go test -v -run=TestLocalStore_simpleStrategy ./...
	