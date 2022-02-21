pre-commit: staticcheck

staticcheck: $(shell go env GOPATH)/bin/staticcheck
	go install honnef.co/go/tools/cmd/staticcheck@latest
