package wallets

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPay(t *testing.T) {
	agent1 := uuid.New()
	agent2 := uuid.New()

	baseWallet := func(addr string, units, quarks int, ownerID uuid.UUID) model.Wallet {
		return model.Wallet{
			Address: addr,
			Units:   units,
			Quarks:  quarks,
			OwnerID: ownerID,
		}
	}
	t.Run("success", func(t *testing.T) {
		expectedFee := int(100*FeeUnits + 10*FeeQuarks)
		from := baseWallet("1", 1000, 10, agent1)
		to := baseWallet("2", 0, 0, agent2)

		newFrom, newTo, fee, cd, err := Pay(from, to, 100, 10)
		require.NoError(t, err)
		assert.Equal(t, expectedFee, fee)
		assert.Equal(t, 100, newTo.Units)
		assert.Equal(t, 10, newTo.Quarks)
		assert.Equal(t, from.Units-100-expectedFee, newFrom.Units)
		assert.Equal(t, 0, newFrom.Quarks)
		assert.Equal(t, time.Second*(110*CooldownPercentage), cd)
	})

	t.Run("insufficient balance (units)", func(t *testing.T) {
		from := baseWallet("1", 100, 10, agent1)
		to := baseWallet("2", 0, 0, agent2)
		_, _, _, _, err := Pay(from, to, 1000, 0)
		require.Error(t, err)
		var withCode errs.WithCode
		if assert.ErrorAs(t, err, &withCode) {
			assert.Equal(t, errs.CodeInsufficientBalance, withCode.Code)
		}
	})

	t.Run("insufficient balance (quarks)", func(t *testing.T) {
		from := baseWallet("1", 100, 10, agent1)
		to := baseWallet("2", 0, 0, agent2)
		_, _, _, _, err := Pay(from, to, 0, 100)
		require.Error(t, err)
		var withCode errs.WithCode
		if assert.ErrorAs(t, err, &withCode) {
			assert.Equal(t, errs.CodeInsufficientBalance, withCode.Code)
		}
	})

	t.Run("insufficient balance (fee)", func(t *testing.T) {
		from := baseWallet("1", 100, 10, agent1)
		to := baseWallet("2", 0, 0, agent2)
		_, _, _, _, err := Pay(from, to, 100, 10)
		require.Error(t, err)
		var withCode errs.WithCode
		if assert.ErrorAs(t, err, &withCode) {
			assert.Equal(t, errs.CodeInsufficientBalance, withCode.Code)
		}
	})

	t.Run("success (no fee)", func(t *testing.T) {
		from := baseWallet("1", 1000, 10, agent1)
		to := baseWallet("2", 0, 0, agent1)

		newFrom, newTo, fee, cd, err := Pay(from, to, 100, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, fee)
		assert.Equal(t, 100, newTo.Units)
		assert.Equal(t, 10, newTo.Quarks)
		assert.Equal(t, 900, newFrom.Units)
		assert.Equal(t, 0, newFrom.Quarks)
		assert.Equal(t, time.Second*(110*CooldownPercentage), cd)
	})
}

func TestGenerateAddress(t *testing.T) {
	addr, err := GenerateAddress()
	require.NoError(t, err)
	assert.Len(t, addr, 43)
	assert.True(t, strings.HasPrefix(addr, "ag1q"))
}

func TestValidateAddress(t *testing.T) {
	addr, _ := GenerateAddress()
	err := ValidateAddress(addr)
	require.NoError(t, err)

	err2 := ValidateAddress("ag1qcju8kvwy3zw6t29m7tdewfu4hhfeh6ckskdpclc")
	require.NoError(t, err2)

	err3 := ValidateAddress("ag1qcju8kvwy3zw6t29m7tdewfu4hhfeh6ckskdpclx")
	require.Error(t, err3)
}
