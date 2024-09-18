binary: css
	go build -o stor ./main.go
css:
	npx tailwindcss -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css -m
css-watch:
	npx tailwindcss -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css --watch
clean:
	rm -f ./stor
	rm -f ./internal/ui/css/style.css