

BUILD_DIR=./bin
SERVICE1=rotator
TESTS=integrations


build:
	go build -o $(BUILD_DIR)/$(SERVICE1) ./cmd/$(SERVICE1)


run:
	./$(BUILD_DIR)/$(SERVICE1)



clean:
	rm -rf $(BUILD_DIR)




up:
	docker-compose up -d --build
	@echo "docker-compose up -d start..."

down:
	docker-compose down
	@echo "docker-compose up -d  stop."


tests:
	go test ./... || (echo "unit test failed"; exit 1)

integration-tests:
	$(MAKE) up
	@echo "цait..."
	@sleep 10 # TODO: нужен ли слип? подвисает
	@echo "bin/integrations-test..."
	go build -o $(BUILD_DIR)/$(TESTS) ./cmd/$(TESTS) && \
    ./$(BUILD_DIR)/$(TESTS) || (echo "tests were failed"; $(MAKE) down; exit 1)
	$(MAKE) down


# TODO: https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
