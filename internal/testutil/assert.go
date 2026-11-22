package testutil

import (
	"testing"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/stretchr/testify/assert"
)

func AssertErrorCode(t *testing.T, err error, errorCode errs.ErrorCode) {
	assert.ErrorIs(t, err, errorCode)
}
