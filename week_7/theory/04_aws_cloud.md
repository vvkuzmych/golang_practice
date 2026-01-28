# AWS Cloud для Go Developers

## AWS SDK Setup

```go
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// Create session
sess := session.Must(session.NewSessionWithOptions(session.Options{
    Config: aws.Config{
        Region: aws.String("us-east-1"),
    },
}))
```

## S3 Operations

```go
// Upload file
func uploadToS3(sess *session.Session, bucket, key string, file io.Reader) error {
    uploader := s3manager.NewUploader(sess)
    
    _, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        Body:   file,
    })
    
    return err
}

// Download file
func downloadFromS3(sess *session.Session, bucket, key string) ([]byte, error) {
    svc := s3.New(sess)
    
    result, err := svc.GetObject(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, err
    }
    defer result.Body.Close()
    
    return io.ReadAll(result.Body)
}

// List objects
func listS3Objects(sess *session.Session, bucket string) error {
    svc := s3.New(sess)
    
    resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        return err
    }
    
    for _, item := range resp.Contents {
        fmt.Printf("Name: %s, Size: %d\n", *item.Key, *item.Size)
    }
    
    return nil
}
```

## Lambda

```go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
    Name string `json:"name"`
}

type MyResponse struct {
    Message string `json:"message"`
}

func HandleRequest(ctx context.Context, event MyEvent) (MyResponse, error) {
    return MyResponse{
        Message: fmt.Sprintf("Hello, %s!", event.Name),
    }, nil
}

func main() {
    lambda.Start(HandleRequest)
}
```

## DynamoDB

```go
import "github.com/aws/aws-sdk-go/service/dynamodb"

// Put item
func putItem(sess *session.Session) error {
    svc := dynamodb.New(sess)
    
    _, err := svc.PutItem(&dynamodb.PutItemInput{
        TableName: aws.String("Users"),
        Item: map[string]*dynamodb.AttributeValue{
            "id":   {S: aws.String("123")},
            "name": {S: aws.String("John")},
        },
    })
    
    return err
}

// Get item
func getItem(sess *session.Session, id string) error {
    svc := dynamodb.New(sess)
    
    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("Users"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {S: aws.String(id)},
        },
    })
    
    return err
}
```

## SQS (Message Queue)

```go
import "github.com/aws/aws-sdk-go/service/sqs"

// Send message
func sendMessage(sess *session.Session, queueURL, message string) error {
    svc := sqs.New(sess)
    
    _, err := svc.SendMessage(&sqs.SendMessageInput{
        QueueUrl:    aws.String(queueURL),
        MessageBody: aws.String(message),
    })
    
    return err
}

// Receive messages
func receiveMessages(sess *session.Session, queueURL string) error {
    svc := sqs.New(sess)
    
    result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(queueURL),
        MaxNumberOfMessages: aws.Int64(10),
        WaitTimeSeconds:     aws.Int64(20), // Long polling
    })
    if err != nil {
        return err
    }
    
    for _, msg := range result.Messages {
        fmt.Printf("Message: %s\n", *msg.Body)
        
        // Delete message after processing
        svc.DeleteMessage(&sqs.DeleteMessageInput{
            QueueUrl:      aws.String(queueURL),
            ReceiptHandle: msg.ReceiptHandle,
        })
    }
    
    return nil
}
```

## Environment Variables

```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
```

## Best Practices

1. Use IAM roles (not hardcoded keys)
2. Enable CloudWatch logging
3. Use KMS for encryption
4. Implement retry logic
5. Handle throttling
6. Use connection pooling
