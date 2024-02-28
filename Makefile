build:
	templ generate
	npx tailwindcss -i ./dist/index.css -o ./dist/tailwind.css
	@go build -o bin/todo_go .

build-dev:
	templ generate
	npx tailwindcss --watch -i ./dist/index.css -o ./dist/tailwind.css
	@go build -o ./tmp/main .

run: build
	./bin/todo_go
