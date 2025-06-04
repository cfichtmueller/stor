binary: css
	go build -o stor ./main.go
css:
	npx tailwindcss@3.4.17 -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css -m
css-watch:
	npx tailwindcss@3.4.17 -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css --watch
clean:
	rm -f ./stor
	rm -f ./internal/ui/css/style.css