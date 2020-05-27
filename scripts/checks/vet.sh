set -eufo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

go vet -printfuncs=Debug,Debugf,Debugln,Info,Infof,Infoln,Notice,Noticef,Noticeln,Error,Errorf,Errorln,Warning,Warningf,Warningln,Critical,Criticalf,Criticalln \
./...
