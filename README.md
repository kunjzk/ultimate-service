# ultimate-service

## Running locally
1. Install Go etc
2. `make run`

### Help
If you need to see help options, edit makefile by adding ` --help` at the end of the `go run` command. Help menu will specify the available config options as well as the relevant command line flags or environmental variables required to edit them.

### Human-readable  stdout logs
Set makefile's `run` command to the following:
```
run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go
```