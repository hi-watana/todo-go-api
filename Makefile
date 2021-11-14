run:
	go run .
test:
	go test . -cover
cover:
	go test . -coverprofile=cover.out; go tool cover -html=cover.out -o cover.html
