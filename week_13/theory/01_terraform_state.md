# Terraform State

## 🎯 Що таке Terraform State?

**Terraform State** - це файл, що зберігає mapping між вашою конфігурацією (`.tf` files) та реальними ресурсами в cloud.

```
Your Config (.tf) ←→ State File ←→ Real Resources (AWS)
```

**Файл:** `terraform.tfstate` (JSON format)

---

## 📊 Чому State важливий?

### Without State

```
You: terraform apply
Terraform: Creates NEW resources every time ❌
Result: Duplicates everywhere!
```

### With State

```
You: terraform apply
Terraform: Checks state → "Resource already exists" ✅
Result: Updates existing resource, no duplicates!
```

---

## 🔍 Що зберігає State?

### State File Structure

```json
{
  "version": 4,
  "terraform_version": "1.6.0",
  "serial": 1,
  "lineage": "unique-id",
  "outputs": {},
  "resources": [
    {
      "mode": "managed",
      "type": "aws_lambda_function",
      "name": "processor",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "arn": "arn:aws:lambda:us-east-1:123456789012:function:processor",
            "function_name": "processor",
            "runtime": "provided.al2023",
            "memory_size": 1024,
            ...
          }
        }
      ]
    }
  ]
}
```

**Зберігає:**
- Resource IDs (ARNs, IDs)
- Resource attributes
- Dependencies між resources
- Metadata (version, serial)

---

## 🏗️ State Backends

### Local Backend (Default)

```hcl
# No configuration needed
# State stored in: ./terraform.tfstate
```

**Pros:**
- ✅ Simple (no setup)
- ✅ Fast (local file)

**Cons:**
- ❌ Not shared (team can't collaborate)
- ❌ Not backed up
- ❌ No locking (concurrent runs corrupt state)
- ❌ Secrets in plain text

**Use for:** Learning, local dev only

### S3 Backend (Recommended)

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"  # For locking
  }
}
```

**Pros:**
- ✅ Shared (team collaboration)
- ✅ Backed up (S3 versioning)
- ✅ Locked (DynamoDB prevents concurrent runs)
- ✅ Encrypted (at rest)

**Cons:**
- ❌ Requires S3 bucket + DynamoDB table
- ❌ Extra AWS costs (minimal)

**Use for:** Production, team environments

---

## 🔐 State Locking

### Problem: Concurrent Modifications

```
Developer 1: terraform apply (starts)
Developer 2: terraform apply (starts)  ← DANGER!

Both modify state simultaneously → Corrupt state! 💥
```

### Solution: DynamoDB Locking

```hcl
terraform {
  backend "s3" {
    dynamodb_table = "terraform-locks"  # Enables locking
  }
}
```

**How it works:**
```
1. Dev 1: terraform apply
2. Terraform acquires lock in DynamoDB
3. Dev 2: terraform apply
4. Terraform checks lock → "Locked by Dev 1" → Waits
5. Dev 1 finishes → Releases lock
6. Dev 2 proceeds
```

**Lock table schema:**
```hcl
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
```

---

## 🔄 State Operations

### Initialize Backend

```bash
terraform init

# Output:
# Initializing the backend...
# Successfully configured the backend "s3"!
```

### View State

```bash
# List resources
terraform state list

# Show resource details
terraform state show aws_lambda_function.processor

# Output entire state
terraform show
```

### Pull State

```bash
# Download state from backend
terraform state pull > state.json

# Useful for inspection
```

### Push State (⚠️ Dangerous)

```bash
# Upload local state to backend
terraform state push terraform.tfstate

# ⚠️ Only if you know what you're doing!
```

---

## 🎯 State Management Commands

### Move Resource

```bash
# Rename in state without recreating
terraform state mv \
  aws_lambda_function.old_name \
  aws_lambda_function.new_name
```

### Remove Resource

```bash
# Remove from state (resource still exists in AWS)
terraform state rm aws_lambda_function.processor

# Terraform will no longer manage this resource
```

### Import Existing Resource

```bash
# Add existing AWS resource to state
terraform import \
  aws_lambda_function.processor \
  processor

# Now Terraform manages this Lambda
```

---

## 🔒 State Security

### Sensitive Data in State

**⚠️ Problem:** State может містити secrets!

```json
{
  "resources": [
    {
      "type": "aws_db_instance",
      "instances": [
        {
          "attributes": {
            "password": "supersecret123"  ← Visible in state!
          }
        }
      ]
    }
  ]
}
```

### Mitigation Strategies

**1. Encrypt State at Rest**
```hcl
terraform {
  backend "s3" {
    encrypt = true  # ✅ Encrypts with AWS KMS
  }
}
```

**2. Restrict Access (S3 Bucket Policy)**
```hcl
resource "aws_s3_bucket_policy" "state" {
  bucket = aws_s3_bucket.terraform_state.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Deny"
        Principal = "*"
        Action = "s3:*"
        Resource = "${aws_s3_bucket.terraform_state.arn}/*"
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}
```

**3. Use Secrets Manager**
```hcl
# Don't store password in config
data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = "db-password"
}

resource "aws_db_instance" "main" {
  password = data.aws_secretsmanager_secret_version.db_password.secret_string
}
```

**4. Enable S3 Versioning**
```hcl
resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }
}
```

---

## 🔄 State Drift

### What is Drift?

**Drift** = Difference між state та real resources

```
State says:    memory_size = 512
Real Lambda:   memory_size = 1024  ← Someone changed manually!

Drift detected! ⚠️
```

### Detect Drift

```bash
# Compare state with reality
terraform plan

# Output:
# aws_lambda_function.processor will be updated in-place
#   ~ memory_size = 1024 -> 512
```

### Fix Drift

**Option 1: Apply (Revert Manual Changes)**
```bash
terraform apply

# Reverts Lambda back to 512 MB
```

**Option 2: Refresh (Update State)**
```bash
terraform refresh

# Updates state to match reality (1024 MB)

# Then update config:
memory_size = 1024  # Accept the change
```

**Option 3: Import (If Resource Not in State)**
```bash
terraform import aws_lambda_function.processor processor
```

---

## ⚠️ State Corruption

### Causes

1. **Concurrent modifications** (no locking)
2. **Interrupted `terraform apply`**
3. **Manual state editing** (don't!)
4. **Incompatible Terraform versions**

### Recovery

**1. Restore from Backup (S3 Versioning)**
```bash
# List versions
aws s3api list-object-versions \
  --bucket my-terraform-state \
  --prefix terraform.tfstate

# Download previous version
aws s3api get-object \
  --bucket my-terraform-state \
  --key terraform.tfstate \
  --version-id VERSION_ID \
  terraform.tfstate.backup

# Test with backup
terraform plan -state=terraform.tfstate.backup
```

**2. Re-import Resources**
```bash
# If state is lost, re-import all resources
terraform import aws_lambda_function.processor processor
terraform import aws_sqs_queue.main queue-url
# ... etc
```

---

## ✅ Best Practices

### 1. Always Use Remote Backend

```hcl
# ✅ GOOD
terraform {
  backend "s3" {
    bucket = "terraform-state"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
  }
}

# ❌ BAD (for teams)
# No backend = local state only
```

### 2. Enable State Locking

```hcl
# ✅ GOOD
terraform {
  backend "s3" {
    dynamodb_table = "terraform-locks"
  }
}
```

### 3. Enable Versioning

```hcl
resource "aws_s3_bucket_versioning" "state" {
  bucket = aws_s3_bucket.terraform_state.id
  
  versioning_configuration {
    status = "Enabled"  # ✅ Keep old versions
  }
}
```

### 4. Encrypt State

```hcl
terraform {
  backend "s3" {
    encrypt = true  # ✅ Always encrypt
  }
}
```

### 5. Separate State per Environment

```
s3://bucket/dev/terraform.tfstate
s3://bucket/staging/terraform.tfstate
s3://bucket/prod/terraform.tfstate
```

### 6. Never Commit State to Git

```bash
# .gitignore
*.tfstate
*.tfstate.*
.terraform/
```

### 7. Regular Backups

```bash
# Automated backup script
terraform state pull > "backups/state-$(date +%Y%m%d-%H%M%S).json"
```

---

## 🎓 Висновок

### Terraform State:

✅ **Maps config to reality** (no duplicates)  
✅ **Stored remotely** (S3 backend)  
✅ **Locked** (DynamoDB prevents conflicts)  
✅ **Encrypted** (security)  
✅ **Versioned** (recovery)  

### Key Points:

1. **State tracks resources** (IDs, attributes)
2. **Remote backend** for teams (S3 + DynamoDB)
3. **Locking prevents** concurrent modifications
4. **Encrypt + version** for security
5. **Separate state** per environment

### Golden Rule:

**"Never manually edit state, never commit to Git!"**

---

## 📖 Далі

- `02_terraform_plan_apply.md` - Plan & Apply workflow
- `practice/01_lambda_terraform/` - Lambda infrastructure
- `practice/02_sqs_terraform/` - SQS infrastructure

**"State is the source of truth!" 🏗️**
