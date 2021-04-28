package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// GravityDenomPrefix indicates the prefix for all assests minted by this module
	GravityDenomPrefix = ModuleName

	// GravityDenomSeparator is the separator for gravity denoms
	GravityDenomSeparator = "/"

	// ETHContractAddressLen is the length of contract address strings
	ETHContractAddressLen = 42

	// GravityDenomLen is the length of the denoms generated by the gravity module
	GravityDenomLen = len(GravityDenomPrefix) + len(GravityDenomSeparator) + ETHContractAddressLen
)

// ValidateGravityDenom validates that the given denomination is either:
//
//  - A valid base denomination (eg: 'uatom')
//  - A valid gravity bridge token representation (i.e 'gravity/{address}')
func ValidateGravityDenom(denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return err
	}

	denomSplit := strings.SplitN(denom, GravityDenomSeparator, 2)

	switch {
	case strings.TrimSpace(denom) == "",
		len(denomSplit) == 1 && denomSplit[0] == GravityDenomPrefix,
		len(denomSplit) == 2 && (denomSplit[0] != GravityDenomPrefix || strings.TrimSpace(denomSplit[1]) == ""):
		return sdkerrors.Wrapf(ErrInvalidGravityDenom, "denomination should be prefixed with the format '%s%s{address}'", GravityDenomPrefix, GravityDenomSeparator)

	case denomSplit[0] == denom && strings.TrimSpace(denom) != "":
		// denom source is from the current chain. Return nil as it has already been validated
		return nil
	}

	// denom source is ethereum. Validate the ethereum hex address
	if err := ValidateEthAddress(denomSplit[1]); err != nil {
		return fmt.Errorf("invalid contract address: %w", err)
	}

	return nil
}

// GravityDenom returns the prefixed coin denomination of the ERC20 token sdk.Coin in the following
// format: gravity-{address}. Example:
// 	gravity/0xa478c2975ab1ea89e8196811f51a7b7ade33eb11
func GravityDenom(contractAddress string) string {
	return fmt.Sprintf("%s%s%s", GravityDenomPrefix, GravityDenomSeparator, contractAddress)
}

func IsEthereumERC20Token(denom string) bool {
	prefix := GravityDenomPrefix + GravityDenomSeparator
	return strings.HasPrefix(denom, prefix)
}

func IsCosmosCoin(denom string) bool {
	return !IsEthereumERC20Token(denom)
}

func GravityDenomToERC20Contract(denom string) string {
	fullPrefix := GravityDenomPrefix + GravityDenomSeparator
	return strings.TrimPrefix(denom, fullPrefix)
}

func (e ERC20ToDenom) Validate() error {
	if err := sdk.ValidateDenom(e.Denom); err != nil {
		return fmt.Errorf("invalid cosmos denomination: %w", err)
	}
	if err := ValidateEthAddress(e.Erc20Address); err != nil {
		return fmt.Errorf("invalid erc20 address: %w", err)
	}
	return nil
}
