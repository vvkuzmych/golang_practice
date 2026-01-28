.PHONY: help init build zip deploy send send-batch send-invalid logs monitor destroy clean test-load

QUEUE_URL = $(shell terraform output -raw queue_url 2>/dev/null)
DLQ_URL = $(shell terraform output -raw dlq_url 2>/dev/null)
LOG_GROUP = $(shell terraform output -raw log_group_name 2>/dev/null)
LAMBDA_NAME = $(shell terraform output -raw lambda_name 2>/dev/null)

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Initialize Go module
	go mod init order-processor || true
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/events
	go mod tidy

build: ## Build Lambda binary
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap main.go
	@echo "✅ Lambda binary built"

zip: build ## Create deployment package
	zip -j function.zip bootstrap
	@echo "✅ function.zip created"

deploy: zip ## Deploy full stack
	terraform init -upgrade
	terraform apply -auto-approve
	@echo ""
	@echo "✅ Full stack deployed!"
	@echo ""
	@echo "Queue URL: $(QUEUE_URL)"
	@echo "Lambda: $(LAMBDA_NAME)"
	@echo ""
	@echo "Try: make send"

send: ## Send single test message
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending order..."
	aws sqs send-message \
		--queue-url $(QUEUE_URL) \
		--message-body '{"order_id":"ORDER-001","amount":99.99,"customer_id":"CUST-123","timestamp":"$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"}'
	@echo "✅ Message sent! Check logs with: make logs"

send-batch: ## Send 10 test messages
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending 10 orders..."
	@for i in {1..10}; do \
		aws sqs send-message \
			--queue-url $(QUEUE_URL) \
			--message-body "{\"order_id\":\"ORDER-$$i\",\"amount\":$$((100 * i)),\"customer_id\":\"CUST-$$i\",\"timestamp\":\"$(shell date -u +%Y-%m-%dT%H:%M:%SZ)\"}" \
			> /dev/null; \
		echo "Sent $$i/10"; \
	done
	@echo "✅ 10 messages sent! Check logs with: make logs"

send-invalid: ## Send invalid message (triggers DLQ)
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending poison message..."
	aws sqs send-message \
		--queue-url $(QUEUE_URL) \
		--message-body 'INVALID_JSON_POISON'
	@echo "✅ Poison message sent!"
	@echo "After 3 failures, it will move to DLQ (check in ~2 min)"
	@echo "Check DLQ with: make check-dlq"

logs: ## Tail Lambda logs
	@if [ -z "$(LOG_GROUP)" ]; then \
		echo "❌ Lambda not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Tailing logs from: $(LOG_GROUP)"
	aws logs tail $(LOG_GROUP) --follow

logs-errors: ## Show only error logs
	@if [ -z "$(LOG_GROUP)" ]; then \
		echo "❌ Lambda not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	aws logs tail $(LOG_GROUP) --follow --filter-pattern "ERROR"

check-dlq: ## Check DLQ for poison messages
	@if [ -z "$(DLQ_URL)" ]; then \
		echo "❌ DLQ not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Checking DLQ: $(DLQ_URL)"
	aws sqs receive-message \
		--queue-url $(DLQ_URL) \
		--max-number-of-messages 10 \
		--attribute-names All

monitor: ## Show queue and Lambda metrics
	@if [ -z "$(LAMBDA_NAME)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "📊 Queue depth:"
	@aws cloudwatch get-metric-statistics \
		--namespace AWS/SQS \
		--metric-name ApproximateNumberOfMessagesVisible \
		--dimensions Name=QueueName,Value=order-processor-queue \
		--start-time $(shell date -u -d '10 minutes ago' +%Y-%m-%dT%H:%M:%S) \
		--end-time $(shell date -u +%Y-%m-%dT%H:%M:%S) \
		--period 60 \
		--statistics Sum \
		--query 'Datapoints[-1].Sum' \
		--output text
	@echo ""
	@echo "📊 Lambda invocations (last 10 min):"
	@aws cloudwatch get-metric-statistics \
		--namespace AWS/Lambda \
		--metric-name Invocations \
		--dimensions Name=FunctionName,Value=$(LAMBDA_NAME) \
		--start-time $(shell date -u -d '10 minutes ago' +%Y-%m-%dT%H:%M:%S) \
		--end-time $(shell date -u +%Y-%m-%dT%H:%M:%S) \
		--period 60 \
		--statistics Sum \
		--query 'sum(Datapoints[].Sum)' \
		--output text
	@echo ""
	@echo "❌ Lambda errors (last 10 min):"
	@aws cloudwatch get-metric-statistics \
		--namespace AWS/Lambda \
		--metric-name Errors \
		--dimensions Name=FunctionName,Value=$(LAMBDA_NAME) \
		--start-time $(shell date -u -d '10 minutes ago' +%Y-%m-%dT%H:%M:%S) \
		--end-time $(shell date -u +%Y-%m-%dT%H:%M:%S) \
		--period 60 \
		--statistics Sum \
		--query 'sum(Datapoints[].Sum)' \
		--output text

invoke: ## Invoke Lambda directly (test without SQS)
	@if [ -z "$(LAMBDA_NAME)" ]; then \
		echo "❌ Lambda not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Invoking Lambda directly..."
	aws lambda invoke \
		--function-name $(LAMBDA_NAME) \
		--payload '{"Records":[{"messageId":"test-1","body":"{\"order_id\":\"TEST-001\",\"amount\":50.00,\"customer_id\":\"CUST-TEST\"}"}]}' \
		response.json
	@cat response.json
	@rm -f response.json

test-load: ## Load test with 100 messages
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending 100 messages..."
	@for i in {1..100}; do \
		aws sqs send-message \
			--queue-url $(QUEUE_URL) \
			--message-body "{\"order_id\":\"LOAD-$$i\",\"amount\":$$((RANDOM % 1000)),\"customer_id\":\"CUST-$$i\"}" \
			> /dev/null; \
		if [ $$((i % 10)) -eq 0 ]; then \
			echo "Sent $$i/100"; \
		fi; \
	done
	@echo "✅ 100 messages sent!"
	@echo "Monitor processing with: make logs"

test-errors: ## Send 50 valid + 50 invalid messages
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Stack not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending 100 messages (50 valid, 50 invalid)..."
	@for i in {1..100}; do \
		if [ $$((i % 2)) -eq 0 ]; then \
			BODY="{\"order_id\":\"$$i\",\"amount\":100}"; \
		else \
			BODY="INVALID_$$i"; \
		fi; \
		aws sqs send-message \
			--queue-url $(QUEUE_URL) \
			--message-body "$$BODY" \
			> /dev/null; \
	done
	@echo "✅ 100 messages sent (50 will fail)!"
	@echo "Check DLQ in 2 minutes: make check-dlq"

outputs: ## Show Terraform outputs
	terraform output

destroy: ## Destroy all resources
	terraform destroy

clean: ## Clean build artifacts and Terraform files
	rm -f bootstrap function.zip response.json
	rm -rf .terraform terraform.tfstate* .terraform.lock.hcl
	@echo "✅ Cleaned"

rebuild: clean zip ## Clean and rebuild
	@echo "✅ Rebuild complete"

redeploy: rebuild deploy ## Rebuild and redeploy
	@echo "✅ Redeployment complete"
