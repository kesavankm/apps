set -x
set -o errexit
set -o nounset
set -euo pipefail

if [ -z "${OS:-}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION:-}" ]; then
    echo "VERSION must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on

# Enable when moved to vendor approach. See 'go help build'
# export GOFLAGS="-mod=vendor"
export GOFLAGS="-mod=mod"

go mod tidy
go fmt ./...
# go build -o ./bin/app ./...
go install                                                      \
    -installsuffix "static"                                     \
    -ldflags "-X $(go list -m)/pkg/version.Version=${VERSION}"  \
    ./...
