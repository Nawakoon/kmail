package account

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
)

func hexToBytes(hexStr string) ([]byte, error) {
	// Convert the hex string to a big integer
	bigInt, success := new(big.Int).SetString(hexStr, 16)
	if !success {
		return nil, fmt.Errorf("failed to convert hex string to big integer")
	}

	// Convert the big integer to bytes
	bytes := bigInt.Bytes()

	return bytes, nil
}

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func ConnectAccount(privateKeyHex string) (*Account, error) {
	// Convert the private key from hex to bytes
	privateKeyBytes, err := hexToBytes(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Generate the ECDSA private key from the bytes
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)

	// Compute the public key points
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.Curve.ScalarBaseMult(privateKeyBytes)

	return &Account{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (a *Account) Sign(data []byte) ([]byte, error) {
	// Sign the data with the private key
	r, s, err := ecdsa.Sign(rand.Reader, a.PrivateKey, data)
	if err != nil {
		return nil, err
	}

	// Encode the signature
	signature, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// get public key address in hex format
func (a *Account) GetAddress() string {
	return PublicKeyToHex(a.PublicKey)
}

func Verify(publicKey *ecdsa.PublicKey, data, signature []byte) bool {
	// Decode the signature
	var decoded struct{ R, S *big.Int }
	_, err := asn1.Unmarshal(signature, &decoded)
	if err != nil {
		return false
	}

	// Verify the signature
	return ecdsa.Verify(publicKey, data, decoded.R, decoded.S)
}

func ExportPublicKeyPEM(publicKey *ecdsa.PublicKey) string {
	// Encode the public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return ""
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM)
}

func HexToPublicKey(hexStr string) (*ecdsa.PublicKey, error) {
	// Convert the hex string back to bytes
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	// Split the bytes back into X and Y coordinates
	xBytes, yBytes := bytes[:len(bytes)/2], bytes[len(bytes)/2:]

	// Convert the bytes to big.Int
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	// Create a new public key
	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return pubKey, nil
}

func PublicKeyToHex(pubKey *ecdsa.PublicKey) string {
	// Concatenate the X and Y coordinates
	concat := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	// Convert the bytes to a hex string
	hexStr := hex.EncodeToString(concat)

	return hexStr
}
