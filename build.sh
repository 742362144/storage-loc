SHELL_FOLDER=$(dirname $(readlink -f "$0"))
cd ${SHELL_FOLDER}

git pull
cd main
go build server.go
go build benchmark.go