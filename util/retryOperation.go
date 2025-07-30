package util

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// RetryOperation implements a retry mechanism with exponential backoff
// with custom number of attempts and initial delay
func RetryOperation(ctx context.Context, logger *zap.Logger, operation func() error, attempts int, delay time.Duration) error {
	var err error
	for i := range attempts {
		if err = operation(); err == nil {
			// Success
			return nil
		}

		// Else retry
		logger.Error("Operation failed, retrying", zap.Int("attempt", i+1), zap.String("error", err.Error()), zap.Duration("delay", delay))
		select {
		case <-ctx.Done():
			return fmt.Errorf("context has cancelled, cancelling operation")
		case <-time.After(delay):
			// Resume loop
		}
		delay *= 2
	}
	return fmt.Errorf("operation failed after %d attempts. Last error: %w", attempts, err)
}
