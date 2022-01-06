version=devel

build:
	CGO_ENABLED=0 go build -ldflags "-X 'github.com/altinn/dotnet-monitor-sidecar-cli/cmd.versionString=${version}'" -o dist/local/dmsctl

build-goreleaser:
	goreleaser build --single-target --snapshot

test:
	go test ./... -v

lint:
	golint ./...

fmt:
	go fmt ./...

cover:
	go test -coverpkg=./... -coverprofile cover.out ./...
	go tool cover -func cover.out

cover-html:
	go test -coverpkg=./... -coverprofile cover.out ./...
	go tool cover -html cover.out

doc:
	go run docs/generate.go

clean:
	rm -rf dist/
	rm -f cover.out