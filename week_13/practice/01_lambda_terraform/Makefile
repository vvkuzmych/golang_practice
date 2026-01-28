.PHONY: help init build zip plan apply deploy test logs destroy clean

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Initialize Go module
	go mod init lambda-example || true
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/events
	go mod tidy

build: ## Build Lambda binary
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go
	@echo "✅ Lambda binary built: bootstrap"

zip: build ## Create Lambda deployment package
	zip -j function.zip bootstrap
	@echo "✅ Lambda package created: function.zip"

tf-init: ## Initialize Terraform
	terraform init

tf-plan: zip ## Run terraform plan
	terraform plan

tf-apply: zip ## Run terraform apply
	terraform apply

deploy: zip ## Build and deploy (init + apply)
	terraform init -upgrade
	terraform apply -auto-approve
	@echo "✅ Lambda deployed!"
	@echo ""
	@echo "Lambda URL:"
	@terraform output -raw lambda_url
	@echo ""

test: ## Test Lambda function
	@LAMBDA_URL=$$(terraform output -raw lambda_url 2>/dev/null); \
	if [ -z "$$LAMBDA_URL" ]; then \
		echo "❌ Lambda not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi; \
	echo "Testing Lambda at: $$LAMBDA_URL"; \
	curl -X POST $$LAMBDA_URL \
		-H "Content-Type: application/json" \
		-d '{"name":"Terraform","message":"Hello from Makefile!"}'
	@echo ""

logs: ## Tail CloudWatch logs
	@LOG_GROUP=$$(terraform output -raw log_group_name 2>/dev/null); \
	if [ -z "$$LOG_GROUP" ]; then \
		echo "❌ Log group not found. Run 'make deploy' first."; \
		exit 1; \
	fi; \
	echo "Tailing logs from: $$LOG_GROUP"; \
	aws logs tail $$LOG_GROUP --follow

outputs: ## Show Terraform outputs
	terraform output

destroy: ## Destroy all resources
	terraform destroy

clean: ## Clean build artifacts
	rm -f bootstrap function.zip
	rm -rf .terraform terraform.tfstate* .terraform.lock.hcl
	@echo "✅ Cleaned build artifacts"

rebuild: clean zip ## Clean and rebuild
	@echo "✅ Rebuild complete"
