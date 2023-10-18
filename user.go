package utils

import (
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"golang.org/x/crypto/ed25519"
)

type Issuer struct {
	IssuerEd25519PrivateKey      ed25519.PrivateKey
	IssuerEd25519PublicKey       ed25519.PublicKey
	IssuerEd25519PublicKeyBase58 string
}

type User struct {
	UserEd25519PrivateKey  ed25519.PrivateKey
	UserEd25519PublicKey   ed25519.PublicKey
	UserPublicKeyBase58    string
	UserAddressBase58Check string
}

func AddIssuer(t provider.T, hlfProxy HlfProxyService, base58Check string) Issuer {
	var issuerFiatEd25519PrivateKey ed25519.PrivateKey
	var issuerFiatEd25519PublicKey ed25519.PublicKey
	var err error
	var issuerEd25519PublicKeyBase58 string

	t.WithNewStep("Generate cryptos for issuer", func(sCtx provider.StepCtx) {
		issuerFiatEd25519PrivateKey, issuerFiatEd25519PublicKey, err = GetPrivateKeyFromBase58Check(base58Check)
		sCtx.Require().NoError(err)
		issuerEd25519PublicKeyBase58 = base58.Encode(issuerFiatEd25519PublicKey)
	})

	t.WithNewStep("Add issuer. Try to add issuer user in acl, issuer may already exist", func(sCtx provider.StepCtx) {
		_, err = hlfProxy.Invoke("acl", "addUser", issuerEd25519PublicKeyBase58, "test", "testuser", "true")
		if err != nil {
			sCtx.Require().Contains(err.Error(), "already exists")
			return
		}
	})

	time.Sleep(BatchTransactionTimeout)

	t.WithNewStep("Check user is created by querying method `checkKeys` of chaincode `acl`", func(sCtx provider.StepCtx) {
		_, err = hlfProxy.Query("acl", "checkKeys", issuerEd25519PublicKeyBase58)
		sCtx.Require().NoError(err)
	})
	return Issuer{issuerFiatEd25519PrivateKey, issuerFiatEd25519PublicKey, issuerEd25519PublicKeyBase58}
}

func AddUser(t provider.T, hlfProxy HlfProxyService) User {
	var userEd25519PrivateKey ed25519.PrivateKey
	var userEd25519PublicKey ed25519.PublicKey
	var err error
	var userPublicKeyBase58 string
	var userAddressBase58Check string

	t.WithNewStep("Generate cryptos for user", func(sCtx provider.StepCtx) {
		userEd25519PrivateKey, userEd25519PublicKey, err = GeneratePrivateAndPublicKey()
		sCtx.Require().NoError(err)
		userPublicKeyBase58 = base58.Encode(userEd25519PublicKey)
		userAddressBase58Check, err = GetAddressByPublicKey(userEd25519PublicKey)
		sCtx.Require().NoError(err)
	})

	t.WithNewStep("Add user by invoking method `addUser` of chaincode `acl` with valid parameters", func(sCtx provider.StepCtx) {
		res, err := hlfProxy.Invoke("acl", "addUser", userPublicKeyBase58, "test", "testuser", "true")
		sCtx.Require().NoError(err)
		sCtx.Require().NotNil(res)
	})

	time.Sleep(BatchTransactionTimeout)

	t.WithNewStep("Check user is created by querying method `checkKeys` of chaincode `acl`", func(sCtx provider.StepCtx) {
		_, err = hlfProxy.Query("acl", "checkKeys", userPublicKeyBase58)
		sCtx.Require().NoError(err)
	})

	return User{userEd25519PrivateKey, userEd25519PublicKey, userPublicKeyBase58, userAddressBase58Check}
}

func GenerateUserPublicKeyBase58(t provider.T) string {
	var userEd25519PublicKey ed25519.PublicKey
	var err error
	var userPublicKeyBase58 string

	t.WithNewStep("Generate cryptos for user", func(sCtx provider.StepCtx) {
		_, userEd25519PublicKey, err = GeneratePrivateAndPublicKey()
		sCtx.Require().NoError(err)
		userPublicKeyBase58 = base58.Encode(userEd25519PublicKey)
		_, err = GetAddressByPublicKey(userEd25519PublicKey)
		sCtx.Require().NoError(err)
	})
	return userPublicKeyBase58
}

func AddUserGetResponce(t provider.T, hlfProxy HlfProxyService) (User, *Response) {
	var (
		userEd25519PrivateKey  ed25519.PrivateKey
		userEd25519PublicKey   ed25519.PublicKey
		err                    error
		userPublicKeyBase58    string
		userAddressBase58Check string
		res                    *Response
	)

	t.WithNewStep("Generate cryptos for user", func(sCtx provider.StepCtx) {
		userEd25519PrivateKey, userEd25519PublicKey, err = GeneratePrivateAndPublicKey()
		sCtx.Require().NoError(err)
		userPublicKeyBase58 = base58.Encode(userEd25519PublicKey)
		userAddressBase58Check, err = GetAddressByPublicKey(userEd25519PublicKey)
		sCtx.Require().NoError(err)
	})

	t.WithNewStep("Add user by invoking method `addUser` of chaincode `acl` with valid parameters", func(sCtx provider.StepCtx) {
		res, err = hlfProxy.Invoke("acl", "addUser", userPublicKeyBase58, "test", "testuser", "true")
		sCtx.Require().NoError(err)
		sCtx.Require().NotNil(res)
	})

	time.Sleep(BatchTransactionTimeout)
	t.WithNewStep("Check user is created by querying method `checkKeys` of chaincode `acl`", func(sCtx provider.StepCtx) {
		_, err = hlfProxy.Query("acl", "checkKeys", userPublicKeyBase58)
		sCtx.Require().NoError(err)
	})

	return User{userEd25519PrivateKey, userEd25519PublicKey, userPublicKeyBase58, userAddressBase58Check}, res
}
