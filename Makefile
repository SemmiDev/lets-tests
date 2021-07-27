run:
	./cmds/env .env go run main.go
build:
	./cmds/env .env go build main.go
removebuild:
	rm ./main
runbuild:
	./cmds/env .env ./main
tests:
	./cmds/env .env go test ./...