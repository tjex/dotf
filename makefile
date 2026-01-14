VERSION := `git describe --tags --match v[0-9]* 2> /dev/null`

build:
	go build -buildvcs=false -ldflags "-X=main.Version=$(VERSION)"

install:
	go install -buildvcs=false -ldflags "-X=main.Version=$(VERSION)"

