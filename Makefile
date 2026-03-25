.PHONY: test-env-up test-env-down test-env-clean test-acc docs

# Start the test environment for manual setup
test-env-up:
	docker compose -f tests/env/docker-compose.yml up -d

# Stop the test environment
test-env-down:
	docker compose -f tests/env/docker-compose.yml down

# Clean up everything (use with caution)
test-env-clean:
	docker compose -f tests/env/docker-compose.yml down -v
	rm -rf tests/env/config/db/*.sqlite

# Save the current DB as the gold template
test-save-gold:
	cp tests/env/config/db/db.sqlite tests/fixtures/db.sqlite.gold

# Reset the active DB from the gold template
test-reset:
	docker compose -f tests/env/docker-compose.yml stop pangolin gerbil
	cp tests/fixtures/db.sqlite.gold tests/env/config/db/db.sqlite
	docker compose -f tests/env/docker-compose.yml start pangolin gerbil
	@echo "Waiting for Pangolin API to be healthy..."
	@until curl -s -f http://localhost:3000/api/v1/ > /dev/null; do sleep 1; done
	@echo "API is up."

# Run Acceptance Tests (requires gold DB and env vars)
test-acc: test-reset
	ASDF_TERRAFORM_VERSION=1.10.0 TF_ACC=1 go test -v ./provider/...

# Generate documentation
docs:
	ASDF_TERRAFORM_VERSION=1.10.0 go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name pangolin

