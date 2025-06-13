# STOR

[CLI](https://github.com/cfichtmueller/storctl) | [go client](https://github.com/cfichtmueller/stor-go-client)

STOR is an object storage released under the MIT license.

The project is very new and still working towards the 1.0.0 release.

## How to build

`make binary`

## How to run

`go run main.go serve`

Then open [http://localhost:8001](http://localhost:8001) to access the console. The buckets are accessible via [http://localhost:8000](http://localhost:8000)

## Configuration

Configuration is done through the environment.

```bash
DATA_DIR=/var/stor     # optional
API_HOST=127.0.0.1     # optional - defaults to empty (bind all)
API_PORT=8000          # optional
CONSOLE_HOST=127.0.0.1 # optional - defaults to empty (bind all)
CONSOLE_PORT=8001      # optional
TRUST_PROXIES=false    # optional - trust X-Forwarded-For headers, defaults to false
```

## Contribute to STOR

This project currently doesn't accept contributions.