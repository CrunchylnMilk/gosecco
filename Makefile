default: lint test

lint:
	golint ./...

test:
	go test -cover -v ./...

deps-dev:
	go get github.com/golang/lint/golint
	go get gopkg.in/check.v1
	go get golang.org/x/tools/cmd/cover
	go get github.com/modocache/gover

deps-dev-u:
	go get -u github.com/golang/lint/golint
	go get -u gopkg.in/check.v1
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/modocache/gover

ci: lint doctor test coveralls

check-go-imports:
	go get golang.org/x/tools/cmd/goimports
	goimports -w .
	git diff --exit-code .

doctor: check-go-imports

# send coverage data to coveralls
coveralls: run-cover
	go get github.com/mattn/goveralls
	goveralls -coverprofile=.coverprofiles/gover.coverprofile -service=travis-ci || true

run-cover: clean-cover
	mkdir -p .coverprofiles
	go test -coverprofile=.coverprofiles/tree.coverprofile     ./tree
	go test -coverprofile=.coverprofiles/constants.coverprofile     ./constants
	go test -coverprofile=.coverprofiles/parser.coverprofile     ./parser
	go test -coverprofile=.coverprofiles/main.coverprofile
	gover .coverprofiles .coverprofiles/gover.coverprofile

clean-cover:
	$(RM) -rf .coverprofiles

# generats an HTML report with coverage information
cover: run-cover
	go tool cover -html=.coverprofiles/gover.coverprofile
