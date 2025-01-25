run: 
	go run ./cmd/web
gen:
	npx rollup -c
	sass ./ui/src/main.scss ./ui/static/gen/main.css
build:
	go build ./cmd/web
	sass ./ui/src/main.scss ./ui/static/gen/main.css
sass:
	sass ./ui/src/main.scss ./ui/static/gen/main.css