set -eufo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v parallel >/dev/null 2>&1 || { echo 'please install parallel or use image that has it'; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"

find "${SCRIPT_DIR}/checks" -type f -name '*.sh' -print0 | \
    parallel -0rv --will-cite --halt soon,fail=1 -j 1 bash {} | \
    sed 's/bash .*\/\(.*\)\.sh/\1/g' \
    || exit 1
