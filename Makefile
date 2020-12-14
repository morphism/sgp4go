all:
	cd cmd/sgp4go && go install

test:
	go test -v -bench=.
