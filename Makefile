

all:
	go build lizard.go


.PHONY: test, bench

test:
	go test ./average
	go test ./statistic


bench:
	go test -test.bench=. ./statistic
