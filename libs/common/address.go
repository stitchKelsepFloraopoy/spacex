package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/dominonetwork/spacex/libs/hexutil"
	"golang.org/x/crypto/sha3"
)

const (
	AddressLength = 20
)

var (
	addressT     = reflect.TypeOf(Address{})
	EmptyAddress = Address{}
)

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if HasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && IsHex(s)
}

//IsLKC check whether a coin is LKC or other token.
func IsLKC(addr Address) bool {
	if addr == EmptyAddress {
		return true
	}
	return false
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + strings.ToLower(string(result))
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, a[:])), nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAddress retains the original string, which may or may not be
// correctly checksummed
type MixedcaseAddress struct {
	addr     Address
	original string
}

// NewMixedcaseAddress constructor (mainly for testing)
func NewMixedcaseAddress(addr Address) MixedcaseAddress {
	return MixedcaseAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcaseAddressFromString is mainly meant for unit-testing
func NewMixedcaseAddressFromString(hexaddr string) (*MixedcaseAddress, error) {
	if !IsHexAddress(hexaddr) {
		return nil, fmt.Errorf("Invalid address")
	}
	a := FromHex(hexaddr)
	return &MixedcaseAddress{addr: BytesToAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcaseAddress
func (ma *MixedcaseAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON(addressT, input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value
func (ma *MixedcaseAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// Address returns the address
func (ma *MixedcaseAddress) Address() Address {
	return ma.addr
}

// String implements fmt.Stringer
func (ma *MixedcaseAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the address has valid checksum
func (ma *MixedcaseAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string
func (ma *MixedcaseAddress) Original() string {
	return ma.original
}
