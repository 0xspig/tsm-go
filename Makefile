run: 
	go run ./cmd/web
gen:
	npx rollup -c
	sass ./ui/src/main.scss ./ui/static/main.css
build:
	go build ./cmd/web
	npx sass ./ui/src/main.scss ./ui/static/main.css
sass:
	npx sass ./ui/src/main.scss ./ui/static/main.css