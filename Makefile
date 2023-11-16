build:
	@go build -o ./bin/geospatial

run: build
	@./bin/geospatial

test:
	@go test ./...