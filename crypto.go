package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

// signMessage - sign arguments with private key in ed25519
func signMessage(privateKey ed25519.PrivateKey, msgToSign []byte) []byte {
	sig := ed25519.Sign(privateKey, msgToSign)
	return sig
}

// verifyEd25519 - verify publicKey with message and signed message
func verifyEd25519(publicKey ed25519.PublicKey, bytesToSign []byte, sMsg []byte) error {
	if !ed25519.Verify(publicKey, bytesToSign, sMsg) {
		err := fmt.Errorf("valid signature rejected")
		return err
	}
	return nil
}

// Sign - sign arguments before send to hlf. create message with certain order arguments expected by chaincode validation in foundation library
func Sign(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, channel string, chaincode string, methodName string, args []string) ([]string, error) {
	return SignWithNonce(privateKey, publicKey, channel, chaincode, methodName, args, "")
}

// SignWithNonce - sign arguments before send to hlf. create message with certain order arguments expected by chaincode validation in foundation library
func SignWithNonce(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, channel string, chaincode string, methodName string, args []string, nonce string) ([]string, error) {
	if nonce == "" {
		nonce = GetNonce()
	}

	return sign(privateKey, publicKey, channel, chaincode, methodName, args, nonce, "")
}

// SignHex - sign arguments in HEX before send to hlf. create message with certain order arguments expected by chaincode validation in foundation library
func SignHex(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, methodName string, args []string) ([]string, error) {
	nonce := GetNonce()
	return SignHexWithNonce(privateKey, publicKey, methodName, args, nonce)
}

// SignHexWithNonce - sign arguments in HEX before send to hlf. create message with certain order arguments expected by chaincode validation in foundation library
func SignHexWithNonce(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, methodName string, args []string, nonce string) ([]string, error) {
	if nonce == "" {
		nonce = GetNonce()
	}

	msg := []string{methodName}
	msg = append(msg, args...)
	msg = append(msg, nonce)
	msg = append(msg, ConvertPublicKeyToBase58(publicKey))
	bytesToSign := sha3.Sum256([]byte(strings.Join(msg, "")))
	sMsg := signMessage(privateKey, bytesToSign[:])

	err := verifyEd25519(publicKey, bytesToSign[:], sMsg)
	if err != nil {
		return nil, err
	}

	return append(msg[1:], hex.EncodeToString(sMsg)), nil
}

// MultisigHex - added multisign in HEX
func MultisigHex(methodName string, args []string, users ...User) ([]string, error) {
	nonce := GetNonce()
	return MultisigHexWithNonce(methodName, args, nonce, users...)
}

// MultisigHexWithNonce - added multisign in HEX with nonce
func MultisigHexWithNonce(methodName string, args []string, nonce string, users ...User) ([]string, error) {
	if nonce == "" {
		nonce = GetNonce()
	}

	msg := []string{methodName}
	msg = append(msg, args...)
	msg = append(msg, nonce)

	for _, i := range users {
		msg = append(msg, i.UserPublicKeyBase58)
	}

	bytesToSign := sha3.Sum256([]byte(strings.Join(msg, "")))

	for _, i := range users {
		sMsg := signMessage(i.UserEd25519PrivateKey, bytesToSign[:])
		msg = append(msg, hex.EncodeToString(sMsg))
		err := verifyEd25519(i.UserEd25519PublicKey, bytesToSign[:], sMsg)
		if err != nil {
			return nil, err
		}
	}
	return msg[1:], nil
}

// GeneratePrivateAndPublicKey - create new private and public key
func GeneratePrivateAndPublicKey() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	return privateKey, publicKey, err
}

// GetAddressByPublicKey - get address by encoded string in standard encoded for project is 'base58.Check'
func GetAddressByPublicKey(publicKey ed25519.PublicKey) (string, error) {
	if len(publicKey) == 0 {
		return "", errors.New("publicKey can't be empty")
	}

	hash := sha3.Sum256(publicKey)
	return base58.CheckEncode(hash[1:], hash[0]), nil
}

// GetPrivateKeyFromBase58Check - get private key type Ed25519 by string - Base58Check encoded private key
func GetPrivateKeyFromBase58Check(secretKey string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	decode, ver, err := base58.CheckDecode(secretKey)
	if err != nil {
		return nil, nil, fmt.Errorf("check decode: %w", err)
	}
	privateKey := ed25519.PrivateKey(append([]byte{ver}, decode...))
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, errors.New("type assertion failed")
	}
	return privateKey, publicKey, nil
}

// ConvertPublicKeyToBase58 - use publicKey with standard encoded type - Base58
func ConvertPublicKeyToBase58(publicKey ed25519.PublicKey) string {
	return base58.Encode(publicKey)
}

// SignExpand - sign arguments before send to hlf. create message with certain order arguments expected by chaincode validation in foundation library
func SignExpand(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, channel string, chaincode string, methodName string, args []string, nonce string, externalRequestID string) ([]string, error) {
	return sign(privateKey, publicKey, channel, chaincode, methodName, args, nonce, externalRequestID)
}

func sign(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, channel string, chaincode string, methodName string, args []string, nonce string, requestID string) ([]string, error) {
	if nonce == "" {
		return nil, errors.New("undefined nonce")
	}

	msg := append(append([]string{methodName, requestID, chaincode, channel}, args...), nonce, ConvertPublicKeyToBase58(publicKey))
	bytesToSign := sha3.Sum256([]byte(strings.Join(msg, "")))
	sMsg := signMessage(privateKey, bytesToSign[:])

	err := verifyEd25519(publicKey, bytesToSign[:], sMsg)
	if err != nil {
		return nil, err
	}

	return append(msg[1:], base58.Encode(sMsg)), nil
}
