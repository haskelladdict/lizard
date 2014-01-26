

all:
	go build lizard.go


.PHONY: test, bench

test:
	go test ./average
	go test ./statistic
	go test ./quickselect


bench:
	go test -test.bench=. ./statistic
