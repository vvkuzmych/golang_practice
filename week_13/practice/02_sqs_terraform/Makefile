.PHONY: help init plan apply deploy send receive delete-msg check-dlq monitor destroy clean

QUEUE_URL = $(shell terraform output -raw queue_url 2>/dev/null)
DLQ_URL = $(shell terraform output -raw dlq_url 2>/dev/null)
QUEUE_NAME = $(shell terraform output -raw queue_name 2>/dev/null)

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Initialize Terraform
	terraform init

plan: ## Run terraform plan
	terraform plan

apply: ## Run terraform apply
	terraform apply

deploy: ## Deploy (init + apply auto-approve)
	terraform init -upgrade
	terraform apply -auto-approve
	@echo "✅ SQS queues deployed!"
	@echo ""
	@echo "Queue URL: $(QUEUE_URL)"
	@echo "DLQ URL: $(DLQ_URL)"

send: ## Send test message
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Sending message to $(QUEUE_URL)"
	aws sqs send-message \
		--queue-url $(QUEUE_URL) \
		--message-body '{"order_id":"123","amount":100,"timestamp":"$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"}'
	@echo "✅ Message sent!"

send-batch: ## Send 10 test messages
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@for i in {1..10}; do \
		echo "Sending message $$i/10..."; \
		aws sqs send-message \
			--queue-url $(QUEUE_URL) \
			--message-body "{\"order_id\":\"$$i\",\"amount\":$$((100 * i))}" \
			> /dev/null; \
	done
	@echo "✅ Sent 10 messages!"

receive: ## Receive messages (long polling)
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Receiving messages from $(QUEUE_URL) (long polling 20s)..."
	aws sqs receive-message \
		--queue-url $(QUEUE_URL) \
		--max-number-of-messages 10 \
		--wait-time-seconds 20 \
		--attribute-names All \
		--message-attribute-names All

delete-msg: ## Delete message (provide RECEIPT_HANDLE)
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@if [ -z "$(RECEIPT_HANDLE)" ]; then \
		echo "❌ RECEIPT_HANDLE not provided. Usage: make delete-msg RECEIPT_HANDLE=..."; \
		exit 1; \
	fi
	aws sqs delete-message \
		--queue-url $(QUEUE_URL) \
		--receipt-handle $(RECEIPT_HANDLE)
	@echo "✅ Message deleted!"

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

purge: ## Purge main queue (delete all messages)
	@if [ -z "$(QUEUE_URL)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "⚠️  Purging queue: $(QUEUE_URL)"
	aws sqs purge-queue --queue-url $(QUEUE_URL)
	@echo "✅ Queue purged!"

purge-dlq: ## Purge DLQ
	@if [ -z "$(DLQ_URL)" ]; then \
		echo "❌ DLQ not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "⚠️  Purging DLQ: $(DLQ_URL)"
	aws sqs purge-queue --queue-url $(DLQ_URL)
	@echo "✅ DLQ purged!"

monitor: ## Show queue metrics
	@if [ -z "$(QUEUE_NAME)" ]; then \
		echo "❌ Queue not deployed. Run 'make deploy' first."; \
		exit 1; \
	fi
	@echo "Queue metrics for: $(QUEUE_NAME)"
	@echo ""
	@echo "📊 Messages visible:"
	aws cloudwatch get-metric-statistics \
		--namespace AWS/SQS \
		--metric-name ApproximateNumberOfMessagesVisible \
		--dimensions Name=QueueName,Value=$(QUEUE_NAME) \
		--start-time $(shell date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
		--end-time $(shell date -u +%Y-%m-%dT%H:%M:%S) \
		--period 300 \
		--statistics Sum \
		--query 'Datapoints[-1].Sum'
	@echo ""
	@echo "📊 Messages in-flight:"
	aws cloudwatch get-metric-statistics \
		--namespace AWS/SQS \
		--metric-name ApproximateNumberOfMessagesNotVisible \
		--dimensions Name=QueueName,Value=$(QUEUE_NAME) \
		--start-time $(shell date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%S) \
		--end-time $(shell date -u +%Y-%m-%dT%H:%M:%S) \
		--period 300 \
		--statistics Sum \
		--query 'Datapoints[-1].Sum'

alarms: ## Show CloudWatch alarms status
	@echo "CloudWatch Alarms:"
	aws cloudwatch describe-alarms \
		--alarm-names \
			my-queue-dlq-not-empty \
			my-queue-high-depth \
			my-queue-old-messages \
		--query 'MetricAlarms[].{Name:AlarmName,State:StateValue}' \
		--output table

outputs: ## Show Terraform outputs
	terraform output

destroy: ## Destroy all resources
	terraform destroy

clean: ## Clean Terraform files
	rm -rf .terraform terraform.tfstate* .terraform.lock.hcl
	@echo "✅ Cleaned Terraform files"

test-poison: ## Send poison message and simulate failure
	@echo "Sending poison message..."
	aws sqs send-message \
		--queue-url $(QUEUE_URL) \
		--message-body "INVALID_JSON_THIS_WILL_FAIL"
	@echo "✅ Poison message sent!"
	@echo ""
	@echo "Now simulate 3 failed processing attempts:"
	@for i in {1..3}; do \
		echo "Attempt $$i/3..."; \
		RECEIPT=$$(aws sqs receive-message \
			--queue-url $(QUEUE_URL) \
			--max-number-of-messages 1 \
			--visibility-timeout 10 \
			--query 'Messages[0].ReceiptHandle' \
			--output text); \
		if [ "$$RECEIPT" != "None" ]; then \
			echo "Received, NOT deleting (simulating failure)"; \
			sleep 11; \
		fi; \
	done
	@echo ""
	@echo "Checking DLQ (wait 10s for message to move)..."
	@sleep 10
	@$(MAKE) check-dlq
