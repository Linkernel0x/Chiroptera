package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateKeyPair Handles the generation of key pairs for developer's authentication
func GenerateKeyPair() ([]byte, []byte, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating key pair: %v", err)
	}
	return pubKey, privKey, nil
}

// SignData signs data using the private key, returns the signature
func SignData(privBytes []byte, data []byte) (string, error) {
	if len(privBytes) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("error signing data: invalid private key")
	}

	signature := ed25519.Sign(privBytes, data)
	return hex.EncodeToString(signature), nil
}

// EvaluateHash calculates the hash of the provided data and returns it
func EvaluateHash(data []byte) (string, error) {
	hashBytes := sha256.Sum256(data)

	hashString := hex.EncodeToString(hashBytes[:])

	return hashString, nil
}

// VerifySignature checks however the signature is valid
func VerifySignature(publicKeyHex string, data []byte, signatureHex string) (bool, error) {
	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("public key not valid: %w", err)
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("signature not valid: %w", err)
	}

	ok := ed25519.Verify(pubBytes, data, sigBytes)
	return ok, nil
}
