# Goroutine Leak Fix

## Problem
The `ShowRefreshing()` method was starting a goroutine that could potentially leak if `HideRefreshing()` was not called or if the application terminated unexpectedly. The goroutine only checked `d.isRefreshing` but had no proper cancellation mechanism.

## Solution
Implemented proper context-based cancellation with the following changes:

### 1. Dashboard Structure Updates
Added context cancellation fields to the Dashboard struct:
```go
refreshingCtx       context.Context
refreshingCancel    context.CancelFunc
```

### 2. ShowRefreshing() Method
- Added proper cancellation of any existing refresh animation before starting a new one
- Created a context with cancel function for each refresh animation
- Used `time.Ticker` instead of `time.Sleep` for better resource management
- Added `select` statement with context cancellation check
- Proper cleanup with `defer ticker.Stop()`

### 3. HideRefreshing() Method
- Calls the cancel function to immediately stop the goroutine
- Cleans up the context and cancel function references
- Maintains the existing UI update functionality

### 4. Cleanup() Method
- Added a public cleanup method to cancel any running goroutines
- Can be called when the dashboard is being destroyed

### 5. Main Application Integration
- Added `defer dashboard.Cleanup()` in the main application to ensure cleanup on exit
- Leverages existing context cancellation pattern in the main app

## Benefits
- **No Goroutine Leaks**: Goroutines are properly cancelled when no longer needed
- **Immediate Cancellation**: Context-based cancellation provides instant termination
- **Resource Efficiency**: Proper cleanup of tickers and goroutines
- **Graceful Shutdown**: Application exit properly cleans up all UI resources
- **Defensive Programming**: Multiple safeguards against resource leaks

## Technical Details
- Uses `context.WithCancel()` for proper cancellation semantics
- Employs `time.Ticker` for better resource management than `time.Sleep`
- Implements defensive programming with fallback checks
- Maintains thread safety with proper synchronization
- Follows Go best practices for goroutine lifecycle management

This fix ensures that the refresh animation goroutines are properly managed and cleaned up, preventing resource leaks and ensuring the application can terminate cleanly under all conditions.