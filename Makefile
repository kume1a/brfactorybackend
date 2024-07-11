build:
	go build

clean:
	rm -f build/

protogen:
	protoc --go_out=. \
	    --go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			./internal/igservice/igservice.proto
run:
	go run main.go serve