BINARY_NAME=goframeworkApp

build:
	@go mod vendor
	@echo "Building Goframework..."
	@go build -o tmp/${BINARY_NAME} .
	@echo  "Goframework built!"

run: build
	@echo "Starting Goframework..."
	@./tmp/${BINARY_NAME} &
	@echo  "Goframework started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

start_compose:
	docker-compose up -d

stop_compose:
	docker-compose down

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping Goframework..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Goframework!"

restart: stop start
