BACKEND_COMMANDS := migrate-up migrate-down migrate-status migrate-create sync-start plan add-emp

$(BACKEND_COMMANDS):
	cd backend && make $@

devdb:
	docker-compose -f docker-compose.dev.yml up postgres -d

.PHONY: $(BACKEND_COMMANDS) devdb