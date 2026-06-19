package nanoid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/uptrace/bun/driver/pgdriver"
)

func Generate(ctx context.Context, tryCreate func(ctx context.Context, extID string) error) (string, error) {
	lengths := []int{4, 6, 8}
	maxAttemptsPerLength := 3

	for _, length := range lengths {
		for attempt := 0; attempt < maxAttemptsPerLength; attempt++ {
			id, err := randomHex(length)
			if err != nil {
				return "", err
			}

			err = tryCreate(ctx, id)
			if err == nil {
				return id, nil // success
			}

			// If it's a duplicate key error, try again
			if isDuplicateKeyError(err) {
				continue
			}
			// Any other error (validation, DB down) → fail immediately
			return "", err
		}
		// fallthrough to next longer length
	}
	return "", errors.New("failed to generate unique ext ID after all attempts")
}

func randomHex(length int) (string, error) {
	bytes := make([]byte, (length+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("randomHex: %w", err)
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr pgdriver.Error
	if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
		return true
	}
	return false
}
