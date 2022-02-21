pre-commit: staticcheck vet fmt

staticcheck: $(shell go env GOPATH)/bin/staticcheck
	staticcheck ./...

$(shell go env GOPATH)/bin/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

vet:
	go vet ./...

fmt:
	go fmt ./...
