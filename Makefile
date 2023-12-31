.PHONY: migrate-create migrate-up migrate-down migrate-force init

PWD = $(shell pwd)
PORT=5432

# Default number of migrations to execute up or down
N = 1

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(PWD)/migrations -seq -digits 5 $(NAME)

migrate-up:
	migrate -database postgres://postgres:password@localhost:$(PORT)/seternak?sslmode=disable -path $(PWD)/migrations up $(N)

migrate-down:
	migrate -database postgres://postgres:password@localhost:$(PORT)/seternak?sslmode=disable -path $(PWD)/migrations down $(N)

migrate-force:
	migrate -database postgres://postgres:password@localhost:$(PORT)/seternak?sslmode=disable -path $(PWD)/migrations force $(VERSION)

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(PWD)/rsa_private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(PWD)/rsa_private.pem -pubout -out $(PWD)/rsa_public.pem

# create dev and test keys
# run postgres containers in docker compose
# migrate down
# migrate up
# docker compose down
init:
	docker compose up -d postgres-seternak && \
	sleep 1
	$(MAKE) migrate-down N= && \
	$(MAKE) migrate-up N= && \
	docker compose down
