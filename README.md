# tars

Twin Pick's Core

## Usage

You can either use [Taskfile](https://taskfile.dev/) or `go` command lines :

```bash
task dev
# or
go run main.go
```

HTTP server will be running on `localhost:8080`.

> [!WARNING]
> In order to use `tars`, you need to ensure that [wall-e](https://github.com/twin-pick/wall-e) is running locally.

## Why Go?

Go is a statically typed, compiled language designed for simplicity and efficiency. It offers excellent performance and strong concurrency support. It seems to be a good fit for building a core service like `tars`, which requires high performance and reliability.

## How it works?

`tars` is basically a HTTP server that handles requests from the application [jarvis](https://github.com/twin-pick/jarvis) (Mobile & Web). It processes the requests and interacts with our Letterboxd Scrapper [wall-e](https://github.com/twin-pick/wall-e) to fetch user's data. Later, it will be responsible for managing the data from various external apis and providing a unified interface for the application.

## License

This project is under [MIT License](LICENSE)
