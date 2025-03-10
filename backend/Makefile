
migrate-up:
	go run migration/cmd/main.go migrate up

migrate-down:
	go run migration/cmd/main.go migrate down

migrate-status:
	go run migration/cmd/main.go migrate status


migrate-create:
	@read -p "Migration açıklaması: " desc; \
	./scripts/create_migration.sh "$$desc"

sync-start:
	go run sync/main.go start

plan:
	go run sync/main.go plan

add-emp:
	go run sync/main.go add-emp

.PHONY: devdb migrate-up migrate-down migrate-status migrate-create sync-start plan
