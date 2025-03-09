devdb:
	docker-compose -f docker-compose.dev.yml up mysql -d

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

.PHONY: devdb migrate-up migrate-down migrate-status migrate-create sync-schema