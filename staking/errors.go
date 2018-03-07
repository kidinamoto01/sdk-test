// nolint
package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

var (
	errCandidateEmpty     = ("Cannot bond to an empty candidate")
	errBadBondingDenom    = ("Invalid coin denomination")
	errBadBondingAmount   = ("Amount must be > 0")
	errNoBondingAcct      = ("No bond account for this (address, validator) pair")
	errCommissionNegative = ("Commission must be positive")
	errCommissionHuge     = ("Commission cannot be more than 100%%")

	errBadValidatorAddr      = ("Validator does not exist for that address")
	errCandidateExistsAddr   = ("Candidate already exist, cannot re-declare candidacy")
	errMissingSignature      = ("Missing signature")
	errBondNotNominated      = ("Cannot bond to non-nominated account")
	errNoCandidateForAddress = ("Validator does not exist for that address")
	errNoDelegatorForAddress = ("Delegator does not contain validator bond")
	errInsufficientFunds     = ("Insufficient bond shares")
	errBadRemoveValidator    = ("Error removing validator")

	invalidInput = sdk.CodeTypeBaseInvalidInput
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code CodeType) string {
	switch code {
	case sdk.CodeInvalidInput:
		return "Invalid input coins"
	case sdk.CodeInvalidOutput:
		return "Invalid output coins"
	case sdk.CodeInvalidAddress:
		return "Invalid address"
	case sdk.CodeUnknownAddress:
		return "Unknown address"
	case sdk.CodeInsufficientCoins:
		return "Insufficient coins"
	case sdk.CodeInvalidCoins:
		return "Invalid coins"
	case sdk.CodeUnknownRequest:
		return "Unknown request"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}

func ErrBadValidatorAddr() sdk.Error {
	return newError(sdk.CodeUnrecognizedAddress, "")
}
func ErrCandidateExistsAddr() sdk.Error {
	return newError(sdk.CodeTypeBaseInvalidInput, errCandidateExistsAddr)
}
func ErrMissingSignature() sdk.Error {
	return newError(sdk.CodeUnauthorized, errMissingSignature)
}
func ErrBondNotNominated() sdk.Error {
	return newError(sdk.CodeInvalidInput, errBondNotNominated)
}
func ErrNoCandidateForAddress() sdk.Error {
	return newError(sdk.CodeUnrecognizedAddress, errNoCandidateForAddress)
}
func ErrNoDelegatorForAddress() sdk.Error {
	return newError(sdk.CodeInvalidInput, errNoDelegatorForAddress)
}
func ErrInsufficientFunds() sdk.Error {
	return newError(sdk.CodeInvalidInput, errInsufficientFunds)
}
func ErrBadRemoveValidator() sdk.Error {
	return newError(sdk.CodeTypeInternalErr, errBadRemoveValidator)
}
