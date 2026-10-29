package wallets

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

const AddressPrefix = "ag1q"
const customAlphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var encoding = base32.NewEncoding(customAlphabet).WithPadding(base32.NoPadding)

// fee = units * FeeUnits + quarks * FeeQuarks
// fee is units!!!
const FeeUnits = 0.01
const FeeQuarks = 5

// cooldown = (units + quarks) * CooldownPercentage * time.Second
const CooldownPercentage = 0.1

// GenerateAddress creates unique wallet address (43 characters)
func GenerateAddress() (string, error) {
	payload := make([]byte, 20)
	if _, err := rand.Read(payload); err != nil {
		return "", err
	}

	checksum := calcChecksum(payload)
	fullBytes := append(payload, checksum...)

	return AddressPrefix + encoding.EncodeToString(fullBytes), nil
}

func ValidateAddress(addr string) error {
	if !strings.HasPrefix(addr, AddressPrefix) {
		return errs.NewWithCode(errs.CodeInvalidWalletAddress, fmt.Errorf("invalid address prefix: %w", errs.ErrInvalidArgument))
	}

	rawPart := strings.TrimPrefix(addr, AddressPrefix)
	decoded, err := encoding.DecodeString(rawPart)
	if err != nil || len(decoded) != 24 {
		return errs.NewWithCode(errs.CodeInvalidWalletAddress, fmt.Errorf("invalid address format: %w", errs.ErrInvalidArgument))
	}

	payload := decoded[:20]
	actualChecksum := decoded[20:]
	expectedChecksum := calcChecksum(payload)

	for i := 0; i < 4; i++ {
		if actualChecksum[i] != expectedChecksum[i] {
			return errs.NewWithCode(errs.CodeInvalidWalletAddress, fmt.Errorf("address checksum mismatch (typo detected): %w", errs.ErrInvalidArgument))
		}
	}
	return nil
}

func calcChecksum(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:4]
}

func Pay(
	fromWallet, toWallet model.Wallet,
	units, quarks int,
) (model.Wallet, model.Wallet, int, time.Duration, error) {
	if fromWallet.Address == toWallet.Address {
		return model.Wallet{}, model.Wallet{}, 0, 0, errs.ErrNotModified
	}

	var fee int
	if fromWallet.OwnerID != toWallet.OwnerID {
		fee = int(float64(units)*FeeUnits + float64(quarks)*FeeQuarks)
	}

	if fromWallet.Units < units+fee {
		return model.Wallet{}, model.Wallet{}, 0, 0, errs.NewWithCode(
			errs.CodeInsufficientBalance,
			fmt.Errorf(
				"insufficient balance (have %d units, need %d): %w",
				fromWallet.Units, units+fee, errs.ErrUnprocessableEntity,
			),
		)
	}

	if fromWallet.Quarks < quarks {
		return model.Wallet{}, model.Wallet{}, 0, 0, errs.NewWithCode(
			errs.CodeInsufficientBalance,
			fmt.Errorf(
				"insufficient balance (have %d quarks, need %d): %w",
				fromWallet.Quarks, quarks, errs.ErrUnprocessableEntity,
			),
		)
	}

	fromWallet.Units -= units + fee
	fromWallet.Quarks -= quarks

	toWallet.Units += units
	toWallet.Quarks += quarks

	cooldown := time.Duration(float64(units+quarks)*CooldownPercentage) * time.Second

	return fromWallet, toWallet, fee, cooldown, nil
}
