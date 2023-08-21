build: 
	@go build -o bin/api

seed: 
	@go run scripts/seed.go

run: build
	@./bin/api	

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 api

test: 
	@go test -v ./...