.PHONY: taiyaku-kun-builder build

taiyaku-kun-builder:
	go build -o taiyaku-kun-builder cmd/taiyaku-kun-builder/main.go

build: taiyaku-kun-builder
	./taiyaku-kun-builder
