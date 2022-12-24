go build -o ./bin/gate cmd/main.go
# services manual https://luci7.medium.com/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-18-04-in-10-minutes-without-docker-e5a1e933bb7e
# see service config in scripts/gate/gate.service
service gate restart
