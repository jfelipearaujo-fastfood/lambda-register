build:
	@echo "Building..."
	@env GOOS=linux GOARCH=arm64 go build -o terraform/bootstrap main.go

zip:
	@echo "Zipping..."
	@zip terraform/lambda.zip terraform/bootstrap

upload:
	@echo "Creating folder..."
	@aws s3api put-object --bucket jsfelipearaujo --key "migrations/" > /dev/null
	@echo "Uploading..."
	@aws s3 cp scripts/ s3://jsfelipearaujo/migrations --recursive

init:
	@echo "Initializing..."
	@cd terraform \
		&& terraform init -reconfigure

check:
	@echo "Checking..."
	make fmt && make validate && make plan

plan:
	@echo "Planning..."
	@cd terraform \
		&& terraform plan -var-file="local.tfvars" -out=plan \
		&& terraform show -json plan > plan.tfgraph

fmt:
	@echo "Formatting..."
	@cd terraform \
		&& terraform fmt -check -recursive

validate:
	@echo "Validating..."
	@cd terraform \
		&& terraform validate

apply:
	@echo "Applying..."
	@cd terraform \
		&& terraform apply plan

destroy:
	@echo "Destroying..."
	@cd terraform \
		&& terraform destroy -auto-approve

gen-tf-docs:
	@echo "Generating Terraform Docs..."
	@terraform-docs markdown table terraform

gen-mocks: ## Gen mock files using mockery
	@if command -v mockery > /dev/null; then \
		echo "Generating..."; \
		mockery; \
	else \
		read -p "Go 'mockery' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/vektra/mockery/v2@latest; \
			echo "Generating..."; \
			mockery; \
		else \
			echo "You chose not to intall mockery. Exiting..."; \
			exit 1; \
		fi; \
	fi

gen-scaffold-bdd:
	@if command -v godog > /dev/null; then \
		echo "Generating BDD scaffold..."; \
		godog ./tests/features; \
	else \
		read -p "Go 'godog' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go go install github.com/cucumber/godog/cmd/godog@latest; \
			echo "Generating BDD scaffold..."; \
			godog ./tests/features; \
		else \
			echo "You chose not to intall godog. Exiting..."; \
			exit 1; \
		fi; \
	fi

test:
	@echo "Running tests..."
	@go test -count=1 ./src/... -v

test-bdd:
	@echo "Running BDD tests..."
	@go test -count=1 ./tests/... -test.v -test.run ^TestFeatures$