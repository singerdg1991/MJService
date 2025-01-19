export ROOT_DIR=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))

convey:
	$(ROOT_DIR)/scripts/convey.sh

# TODO: docker convoy

docker-test:
	chmod +x $(ROOT_DIR)/scripts/test.sh
	$(ROOT_DIR)/scripts/test.sh docker

docker-lint:
	chmod +x $(ROOT_DIR)/scripts/lint.sh
	$(ROOT_DIR)/scripts/lint.sh docker

test:
	chmod +x $(ROOT_DIR)/scripts/test.sh
	$(ROOT_DIR)/scripts/test.sh

lint:
	chmod +x $(ROOT_DIR)/scripts/lint.sh
	$(ROOT_DIR)/scripts/lint.sh

testcov:
	go test -v ./... -coverprofile=$(ROOT_DIR)/sandbox/coverage.out && go tool cover -html=$(ROOT_DIR)/sandbox/coverage.out

cli:
	$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	go run $(ROOT_DIR)/cmd/cli/main.go $(ARGS)

docker-cli:
	$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	docker exec dev_maja_service go run /app/cmd/cli/main.go $(ARGS)

swagger:
	$(ROOT_DIR)/scripts/openapi.sh

# -------------- Development ----------------
start:
	cat $(ROOT_DIR)/.env.development > $(ROOT_DIR)/.env
	chmod +x $(ROOT_DIR)/scripts/start.sh
	$(ROOT_DIR)/scripts/start.sh

docker-start:
	cat $(ROOT_DIR)/.env.development > $(ROOT_DIR)/.env
	chmod +x $(ROOT_DIR)/scripts/start.sh
	$(ROOT_DIR)/scripts/start.sh docker

docker-stop:
	chmod +x $(ROOT_DIR)/scripts/stop.sh
	$(ROOT_DIR)/scripts/stop.sh

# -------------- Production ----------------
serve:
	cat $(ROOT_DIR)/.env.production > $(ROOT_DIR)/.env
	chmod +x $(ROOT_DIR)/scripts/serve.sh
	$(ROOT_DIR)/scripts/serve.sh

docker-serve:
	cat $(ROOT_DIR)/.env.production > $(ROOT_DIR)/.env
	chmod +x $(ROOT_DIR)/scripts/serve.sh
	$(ROOT_DIR)/scripts/serve.sh docker

docker-drop:
	chmod +x $(ROOT_DIR)/scripts/drop.sh
	$(ROOT_DIR)/scripts/drop.sh

# ------------------ Docker -----------------
docker-logs:
	chmod +x $(ROOT_DIR)/scripts/logs.sh
	$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	$(ROOT_DIR)/scripts/logs.sh $(ARGS)

# ------------------- Git -------------------
prepare-hooks:
	chmod +x $(ROOT_DIR)/scripts/prepare-hooks.sh
	$(ROOT_DIR)/scripts/prepare-hooks.sh

gitlog:
	git log --graph --abbrev-commit --decorate --all --pretty=format:'%C(auto)%h%C(reset) %C(bold blue)%an%C(reset) %C(bold green)%ad%C(reset) %C(auto)%C(yellow)%ar%C(reset) %s ' --since='1 year ago'

# ------------------ GRPC -------------------
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./eventstore/protobuf/eventstore.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./tikka/protobuf/tikka.proto