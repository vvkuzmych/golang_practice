# Technical English Communication

## Code Review Vocabulary

### Positive Feedback
- "LGTM" (Looks Good To Me)
- "Nice catch!"
- "Great refactoring"
- "This is elegant"
- "Well tested"

### Constructive Feedback
- "Consider using..."
- "What if we..."
- "Have you thought about..."
- "Could you clarify..."
- "Minor: ..."
- "Nit: ..." (minor issue)

### Issues
- "This breaks..."
- "Edge case: ..."
- "Potential race condition"
- "Memory leak here"
- "Performance concern"

## Daily Standup

**Template:**
- Yesterday: "I completed X, Y"
- Today: "I will work on Z"
- Blockers: "Waiting for..." / "None"

**Example:**
```
Yesterday I finished the user authentication API and wrote unit tests.
Today I'm working on integrating Redis caching.
No blockers.
```

## Email/Slack Communication

### Task Update
```
Hi team,

I've deployed the payment service to staging.
Changes:
- Added Stripe webhook support
- Implemented retry logic
- Updated tests

Ready for QA review.

Staging URL: https://staging.api.example.com
```

### Bug Report
```
Subject: [BUG] Memory leak in worker pool

Description:
The worker pool doesn't properly close channels, causing goroutine leaks.

Steps to reproduce:
1. Start 100 workers
2. Send 10k jobs
3. Monitor memory usage

Expected: Memory stable
Actual: Memory grows continuously

Logs attached. I'll investigate and submit a fix today.
```

### Pull Request Description
```
## Summary
Implement Redis caching for user sessions

## Changes
- Added Redis client wrapper
- Implemented cache-aside pattern
- Added fallback to DB on cache miss
- Added metrics for cache hit/miss

## Testing
- Unit tests for cache operations
- Integration tests with test containers
- Load tested with 10k concurrent users

## Deployment Notes
Requires REDIS_URL environment variable
```

## Technical Documentation

### Function Documentation
```go
// GetUserByID retrieves a user by their unique ID.
// It first checks the cache, then falls back to the database.
//
// Returns ErrUserNotFound if the user doesn't exist.
// Returns ErrDatabaseUnavailable on database errors.
func GetUserByID(ctx context.Context, id int64) (*User, error) {
    // ...
}
```

## Common Phrases

- "Let's schedule a sync"
- "Can we hop on a call?"
- "I'll circle back"
- "Pushing this to next sprint"
- "Taking this offline"
- "FYI" (For Your Information)
- "ETA" (Estimated Time of Arrival)
- "WIP" (Work In Progress)
- "RFC" (Request For Comments)
