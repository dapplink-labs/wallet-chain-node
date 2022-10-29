package keys

import (
	"crypto"
	"crypto/ed25519"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"

	"github.com/mr-tron/base58/base58"
)

// KeyType is the type of key.
type KeyType int

const (
	// ED25519 represents an ed25519 key.
	ED25519 KeyType = iota
)

// NewPublicKeyFromString creates a new public key from a base58 encoded string prefixed with the key type string.
func NewPublicKeyFromString(str string) (*PublicKey, error) {
	parts := strings.Split(str, ":")
	if len(parts) == 1 {
		data, err := base58.Decode(parts[0])
		if err != nil {
			return nil, fmt.Errorf("decoding key string: %v", err)
		}
		return &PublicKey{Type: ED25519, Data: data}, nil
	} else if len(parts) == 2 {
		keyType, err := stringToKeyType(parts[0])
		if err != nil {
			return nil, fmt.Errorf("decoding key type string: %v", err)
		}
		data, err := base58.Decode(parts[1])
		if err != nil {
			return nil, fmt.Errorf("decoding key data: %v", err)
		}
		return &PublicKey{Type: keyType, Data: data}, nil
	} else {
		return nil, fmt.Errorf("invalid encoded key format, must be <curve>:<encoded key>")
	}
}

// PublicKey represents a public key.
type PublicKey struct {
	Type KeyType
	Data []byte
}

// ToString creates a string representation of the public key of the form <curve>:<base58 encoded key>.
func (pk *PublicKey) ToString() (string, error) {
	typeStr, err := keyTypeToString(pk.Type)
	if err != nil {
		return "", fmt.Errorf("converting key type to string: %v", err)
	}
	return fmt.Sprintf("%s:%s", typeStr, base58.Encode(pk.Data)), nil
}

func keyTypeToString(keyType KeyType) (string, error) {
	switch keyType {
	case ED25519:
		return "ed25519", nil
	default:
		return "", fmt.Errorf("unknown key type: %v", keyType)
	}
}

func stringToKeyType(str string) (KeyType, error) {
	switch strings.ToLower(str) {
	case "ed25519":
		return ED25519, nil
	default:
		return -1, fmt.Errorf("unknown key type string: %s", str)
	}
}

// KeyPair represents a public/private key pair.
type KeyPair interface {
	fmt.Stringer
	Sign(message []byte) ([]byte, error)
	Verify(message, signature []byte) bool
	GetPublicKey() PublicKey
}

// NewKeyPairFromRandom creates a random KeyPair using the specified curve.
func NewKeyPairFromRandom(curve string) (KeyPair, error) {
	switch strings.ToUpper(curve) {
	case "ED25519":
		_, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			return nil, fmt.Errorf("generating random ed25519 key: %v", err)
		}
		return &KeyPairEd25519{privateKey: priv}, nil
	default:
		return nil, fmt.Errorf("unknown curve %s", curve)
	}
}

// NewKeyPairFromString creates a new KeyPair from a optionally curve-prefixed base58 string.
func NewKeyPairFromString(secretKey string) (KeyPair, error) {
	parts := strings.Split(secretKey, ":")
	base58string := ""
	if len(parts) == 1 {
		base58string = parts[0]
	} else if len(parts) == 2 {
		switch strings.ToUpper(parts[0]) {
		case "ED25519":
			base58string = parts[1]
		default:
			return nil, fmt.Errorf("unknown curve %s", parts[0])
		}
	} else {
		return nil, fmt.Errorf("Invalid encoded key format, must be <curve>:<encoded key>")
	}
	kp, err := keyPairEd25519FromString(base58string)
	if err != nil {
		return nil, fmt.Errorf("creating ed25519 key from string: %v", err)
	}
	return kp, nil
}

// KeyPairEd25519 is an ed25519 implementation of KeyPair.
type KeyPairEd25519 struct {
	privateKey ed25519.PrivateKey
}

// Sign signs a message with the KeyPair's private key.
func (k *KeyPairEd25519) Sign(message []byte) ([]byte, error) {
	res, err := k.privateKey.Sign(nil, message, crypto.Hash(0))
	if err != nil {
		return nil, fmt.Errorf("calling sign: %v", err)
	}
	return res, nil
}

// Verify reports whether signature is a valid signature of message by the KeyPair's public key.
func (k *KeyPairEd25519) Verify(message, signature []byte) bool {
	return ed25519.Verify(k.privateKey.Public().(ed25519.PublicKey), message, signature)
}

// GetPublicKey returns the PublicKey corresponding to the KeyPair's private key.
func (k *KeyPairEd25519) GetPublicKey() PublicKey {
	return PublicKey{
		Type: ED25519,
		Data: k.privateKey.Public().(ed25519.PublicKey),
	}
}

func (k *KeyPairEd25519) String() string {
	return string(k.privateKey)
}

func keyPairEd25519FromString(base58string string) (*KeyPairEd25519, error) {
	data, err := base58.Decode(base58string)
	if err != nil {
		return nil, fmt.Errorf("decoding secret key: %v", err)
	}
	switch len(data) {
	case ed25519.PrivateKeySize + ed25519.PublicKeySize:
		// Remove the redundant public key. See issue #36.
		redundantPk := data[ed25519.PrivateKeySize:]
		pk := data[ed25519.PrivateKeySize-ed25519.PublicKeySize : ed25519.PrivateKeySize]
		if subtle.ConstantTimeCompare(pk, redundantPk) == 0 {
			return nil, errors.New("expected redundant ed25519 public key to be redundant")
		}

		// No point in storing the extra data.
		newKey := make([]byte, ed25519.PrivateKeySize)
		copy(newKey, data[:ed25519.PrivateKeySize])
		data = newKey
	case ed25519.PrivateKeySize:
	default:
		return nil, fmt.Errorf(
			"expected ed25519 data size to be %d or %d, got %d",
			ed25519.PrivateKeySize,
			ed25519.PrivateKeySize+ed25519.PublicKeySize,
			len(data),
		)
	}
	return &KeyPairEd25519{
		privateKey: ed25519.PrivateKey(data),
	}, nil
}
