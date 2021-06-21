include Makefile.inc
test: get-gpu-setup
	go test ./gpu -v
