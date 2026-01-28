# Security & Compliance

## OWASP Top 10

### 1. Broken Access Control
```go
// ❌ BAD
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    db.Delete(id) // Any user can delete any user!
}

// ✅ GOOD
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    currentUser := getUserFromToken(r)
    id := r.URL.Query().Get("id")
    
    if currentUser.Role != "admin" && currentUser.ID != id {
        http.Error(w, "Forbidden", 403)
        return
    }
    
    db.Delete(id)
}
```

### 2. Injection (SQL Injection)
```go
// ❌ BAD
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
db.Query(query) // Vulnerable!

// ✅ GOOD
db.Query("SELECT * FROM users WHERE email = $1", email)
```

### 3. Sensitive Data Exposure
```go
// ✅ GOOD - Hash passwords
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## JWT Authentication

```go
import "github.com/golang-jwt/jwt/v5"

func GenerateJWT(userID int64) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
}
```

## HIPAA Compliance Basics

### Protected Health Information (PHI)
- Names, addresses, dates
- Medical records
- Social Security numbers
- Health insurance info

### Requirements
✅ Encrypt data at rest and in transit
✅ Implement access controls
✅ Audit logs (who accessed what, when)
✅ Business Associate Agreements (BAA)
✅ Incident response plan

### Implementation
```go
// Encrypt sensitive data
import "crypto/aes"

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    // ... GCM mode encryption
}

// Audit logging
func LogAccess(userID int64, resource string, action string) {
    log.Info().
        Int64("user_id", userID).
        Str("resource", resource).
        Str("action", action).
        Time("timestamp", time.Now()).
        Msg("access_log")
}
```

## GDPR / Data Protection

### Key Principles
- **Right to access** - users can request their data
- **Right to be forgotten** - users can delete their data
- **Data portability** - export user data
- **Consent** - explicit opt-in

### Implementation
```go
// Export user data
func ExportUserData(userID int64) (*UserExport, error) {
    user := db.GetUser(userID)
    orders := db.GetOrders(userID)
    // ... collect all data
    
    return &UserExport{
        User:   user,
        Orders: orders,
    }, nil
}

// Delete user data (GDPR Right to be forgotten)
func DeleteUserData(userID int64) error {
    tx := db.Begin()
    defer tx.Rollback()
    
    // Delete or anonymize all user data
    tx.Delete("orders WHERE user_id = ?", userID)
    tx.Delete("users WHERE id = ?", userID)
    
    return tx.Commit()
}
```

## Secrets Management

### HashiCorp Vault
```go
import "github.com/hashicorp/vault/api"

func getSecret(key string) (string, error) {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return "", err
    }
    
    secret, err := client.Logical().Read("secret/data/" + key)
    if err != nil {
        return "", err
    }
    
    return secret.Data["value"].(string), nil
}
```

### AWS Secrets Manager
```go
import "github.com/aws/aws-sdk-go/service/secretsmanager"

func getAWSSecret(secretName string) (string, error) {
    svc := secretsmanager.New(session.Must(session.NewSession()))
    
    result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    })
    if err != nil {
        return "", err
    }
    
    return *result.SecretString, nil
}
```

## Best Practices

✅ **Never commit secrets** to Git
✅ **Use environment variables** for config
✅ **Encrypt sensitive data** at rest and in transit
✅ **Implement rate limiting**
✅ **Validate all inputs**
✅ **Use HTTPS** everywhere
✅ **Keep dependencies updated**
✅ **Implement proper logging** (but don't log secrets!)
✅ **Regular security audits**
✅ **Principle of least privilege**

---

**Security is not a feature - it's a requirement!** 🔒
