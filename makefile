all: run

.PHONY: depend format build lint clean

depend:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/dep/cmd/dep
	dep ensure -v

AllGOFiles := cmd/*.go pkg/*/*.go

format:
	gofmt -w -s ${AllGOFiles}
	goimports -w ${AllGOFiles}

run: format
	go run cmd/main.go -kubeconfig=${HOME}/.kube/config

