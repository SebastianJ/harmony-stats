__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source ${__dir}/deps.sh
go mod tidy
go fmt ./...
make clean
make all
