run: 
	go run ./cmd/web
gen:
	npx rollup -c
build:
	go build ./cmd/web