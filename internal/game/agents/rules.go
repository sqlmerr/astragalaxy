package agents

import (
	"fmt"
	"regexp"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func checkUsername(username string) error {
	usernameLen := len([]rune(username))
	if usernameLen < 3 || usernameLen > 32 {
		return fmt.Errorf("invalid username length: must be between 3 and 27 characters: %w", errs.CodeInvalidUsername)
	}

	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("invalid username format: must contain only latin letters, digits, and underscores: %w", errs.CodeInvalidUsername)
	}

	return nil
}
