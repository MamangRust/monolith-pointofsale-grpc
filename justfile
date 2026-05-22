COMPOSE_FILE := "deployments/local/docker-compose.yml"
SERVICES := "apigateway migrate auth role user category cashier merchant order_item order product transaction email"
DOCKER_COMPOSE := "docker compose"

migrate:
	go run service/migrate/main.go up

migrate-down:
	go run service/migrate/main.go down

seeder:
	go run service/seeder/main.go


generate-proto:
	cd proto && find . -name "*.proto" -exec protoc --proto_path=. --go_out=../pb --go_opt=paths=source_relative --go-grpc_out=../pb --go-grpc_opt=paths=source_relative {} +


# Build all services that contain a go.mod file
build:
    @mkdir -p bin
    @for mod in service/*/go.mod; do \
        dir=$(dirname $mod); \
        service=$(basename $dir); \
        echo "🔨 Building $service..."; \
        if [ -f "$dir/cmd/main.go" ]; then \
            (cd $dir && go build -o ../../bin/$service ./cmd/main.go) || exit 1; \
        else \
            (cd $dir && go build -o ../../bin/$service ./main.go) || exit 1; \
        fi; \
    done
    @echo "✅ All services built successfully in bin/ folder."


generate-sql:
	sqlc generate


generate-swagger:
	swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs


build-image:
	@for service in $(SERVICES); do \
		echo "🔨 Building $$service-pointofsale-service..."; \
		docker build -t $$service-pointofsale-service:1.0 -f service/$$service/Dockerfile service/$$service || exit 1; \
	done
	@echo "✅ All services built successfully."

image-load:
	@for service in $(SERVICES); do \
		echo "🚚 Loading $$service-service..."; \
		minikube image load $$service-service:1.0 || exit 1; \
	done
	@echo "✅ All services loaded successfully."


ps:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) ps

up:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) up -d

down:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) down

build-up:
	make build-image && make up

# Tidy all go.mod files
tidy-all:
	@for mod in service/*/go.mod; do \
		dir=$(dirname $mod); \
		service=$(basename $dir); \
		echo "🧹 Tidying $service..."; \
		(cd $dir && go mod tidy) || exit 1; \
	done
	@echo "✅ All services tidied successfully."