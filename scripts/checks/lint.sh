set -eufo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v golint >/dev/null 2>&1 || { echo 'please install goimports or use image that has it'; exit 1; }

golint ./...
