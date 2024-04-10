package utils

import (
	"strconv"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// EmitGetTxIDAndCheckBalance emits amount of tokens to userAddressBase58Check and checks that balance is equal to amount
func EmitGetTxIDAndCheckBalance(
	t provider.T,
	hlfProxy HlfProxyService,
	userAddressBase58Check string,
	issuer Issuer,
	channel string,
	chaincode string,
	amount string,
) string {
	var txID string
	t.WithNewStep("Emit "+amount+" token to user "+userAddressBase58Check+" and get txId", func(sCtx provider.StepCtx) {
		emitArgs := []string{userAddressBase58Check, amount}
		signedEmitArgs, err := Sign(issuer.IssuerEd25519PrivateKey, issuer.IssuerEd25519PublicKey, channel, chaincode, "emit", emitArgs)
		sCtx.Require().NoError(err)
		res, err := hlfProxy.Invoke(channel, "emit", signedEmitArgs...)
		sCtx.Require().NoError(err)
		time.Sleep(BatchTransactionTimeout)
		txID = res.TransactionID
	})

	CheckBalanceEqual(t, hlfProxy, userAddressBase58Check, channel, amount)
	return txID
}

// EmitGetResponseAndCheckBalance emits amount of tokens to userAddressBase58Check and checks that balance is equal to amount
func EmitGetResponseAndCheckBalance(
	t provider.T,
	hlfProxy HlfProxyService,
	userAddressBase58Check string,
	iss Issuer,
	channel string,
	chaincode string,
	amount string,
) *Response {
	var res *Response
	t.WithNewStep("Emit "+amount+" token to user "+userAddressBase58Check+" and get txId", func(sCtx provider.StepCtx) {
		emitArgs := []string{userAddressBase58Check, amount}
		signedEmitArgs, err := Sign(iss.IssuerEd25519PrivateKey, iss.IssuerEd25519PublicKey, channel, chaincode, "emit", emitArgs)
		sCtx.Require().NoError(err)
		res, err = hlfProxy.Invoke(channel, "emit", signedEmitArgs...)
		sCtx.Require().NoError(err)
		time.Sleep(BatchTransactionTimeout)
	})

	CheckBalanceEqual(t, hlfProxy, userAddressBase58Check, channel, amount)
	return res
}

// CheckBalanceEqual checks that balance of userAddressBase58Check is equal to amount
func CheckBalanceEqual(t provider.T, hlfProxy HlfProxyService, userAddressBase58Check string, channel string, amount string) {
	t.WithNewStep("Checking that balance equal "+amount, func(sCtx provider.StepCtx) {
		respGetBalance, err := hlfProxy.Query(channel, "balanceOf", userAddressBase58Check)
		sCtx.Require().NoError(err)
		sCtx.Require().Equal("\""+amount+"\"", string(respGetBalance.Payload))
	})
}

// CheckAllowedBalanceEqual checks that allowed balance of userAddressBase58Check is equal to amount
func CheckAllowedBalanceEqual(t provider.T, hlfProxy HlfProxyService, userAddressBase58Check string, channel string, tokenUppercase string, amount string) {
	t.WithNewStep("Checking that allowed balance equal "+amount, func(sCtx provider.StepCtx) {
		resp, err := hlfProxy.Query(channel, "allowedBalanceOf", userAddressBase58Check, tokenUppercase)
		sCtx.Require().NoError(err)
		sCtx.Require().Equal("\""+amount+"\"", string(resp.Payload))
	})
}

// CheckBalanceEqualWithRetry checks that balance of userAddressBase58Check is equal to amount with retries
func CheckBalanceEqualWithRetry(t provider.T, hlfProxy HlfProxyService, userAddressBase58Check string, channel string, amount string, sleep time.Duration, retries int) {
	t.WithNewStep("Checking that balance equal "+amount+" with retry", func(sCtx provider.StepCtx) {
		i := 0
		for i < retries {
			respGetBalance, err := hlfProxy.Query(channel, "balanceOf", userAddressBase58Check)
			t.Require().NoError(err)
			if string(respGetBalance.Payload) == "\""+amount+"\"" {
				return
			}
			t.LogStep("attempt number is " + strconv.Itoa(i+1))
			i++
			time.Sleep(sleep)
		}
		t.LogStep("failed to get expected balance " + amount)
		t.FailNow()
	})
}

// TransferCheckBalanceAndGetRespose transfers amount of tokens from userFrom to userToAddress and checks that balance of userToAddress is equal to amount
func TransferCheckBalanceAndGetRespose(
	t provider.T,
	hlfProxy HlfProxyService,
	userFrom User,
	userToAddress string,
	channel string,
	chaincode string,
	amount string,
) *Response {
	var (
		signedTransferArgs []string
		err                error
		resTransfer        *Response
	)

	t.WithNewStep("Sign arguments before emission process", func(sCtx provider.StepCtx) {
		ref := "ref transfer"
		transferArgs := []string{userToAddress, amount, ref}
		signedTransferArgs, err = Sign(userFrom.UserEd25519PrivateKey, userFrom.UserEd25519PublicKey, channel, chaincode, "transfer", transferArgs)
		sCtx.Require().NoError(err)
	})

	t.WithNewStep("Invoke fiat chaincode by issuer for token emission", func(sCtx provider.StepCtx) {
		resTransfer, err = hlfProxy.Invoke(channel, "transfer", signedTransferArgs...)
		sCtx.Require().NoError(err)
		time.Sleep(BatchTransactionTimeout)
	})

	CheckBalanceEqual(t, hlfProxy, userToAddress, channel, amount)
	return resTransfer
}

// GetEmitPayload emits amount of tokens to userAddressBase58Check and checks that balance is equal to amount
func GetEmitPayload(
	t provider.T,
	hlfProxy HlfProxyService,
	userAddressBase58Check string,
	issuer Issuer,
	amount string,
) string {
	var txID string
	t.WithNewStep("Emit "+amount+" token to user "+userAddressBase58Check+" and get txId", func(sCtx provider.StepCtx) {
		emitArgs := []string{userAddressBase58Check, amount}
		signedEmitArgs, err := Sign(issuer.IssuerEd25519PrivateKey, issuer.IssuerEd25519PublicKey, "inv", "inv", "emit", emitArgs)
		sCtx.Require().NoError(err)
		res, err := hlfProxy.Invoke("fiat", "emit", signedEmitArgs...)
		sCtx.Require().NoError(err)
		time.Sleep(BatchTransactionTimeout)
		txID = res.TransactionID
	})

	CheckBalanceEqual(t, hlfProxy, userAddressBase58Check, "fiat", amount)
	return txID
}

// SwapFiatToCCCheckBalanceAndGetSwapDoneAndSwapBeginTxID swaps amount of tokens from fiat to cc channel
func SwapFiatToCCCheckBalanceAndGetSwapDoneAndSwapBeginTxID(t provider.T, hlfProxy HlfProxyService, user User, amount string) (string, string) {
	var (
		swapBeginTxID string
		swapDoneTxID  string
	)
	t.WithNewStep("Swap between channels", func(sCtx provider.StepCtx) {
		swapBeginArgs := []string{"FIAT", "CC", amount, DefaultSwapHash}
		signedSwapBeginArgs, err := Sign(user.UserEd25519PrivateKey, user.UserEd25519PublicKey, "fiat", "fiat", "swapBegin", swapBeginArgs)
		sCtx.Assert().NoError(err)
		swapBeginResp, err := hlfProxy.Invoke("fiat", "swapBegin", signedSwapBeginArgs...)
		swapBeginTxID = swapBeginResp.TransactionID
		sCtx.Assert().NoError(err)
		time.Sleep(BatchTransactionTimeout)

		sCtx.NewStep("swapGet txID in fiat channel")
		_, err = hlfProxy.Query("fiat", "swapGet", swapBeginResp.TransactionID)
		sCtx.Assert().NoError(err)
		sCtx.NewStep("swapGet txID in cc channel")
		_, err = hlfProxy.Query("cc", "swapGet", swapBeginResp.TransactionID)
		sCtx.Assert().NoError(err)
		sCtx.NewStep("swapDone")
		swapDoneResp, err := hlfProxy.Invoke("cc", "swapDone", swapBeginResp.TransactionID, DefaultSwapKey)
		swapDoneTxID = swapDoneResp.TransactionID
		sCtx.Assert().NoError(err)
		time.Sleep(BatchTransactionTimeout)
		sCtx.NewStep("Get allowed balance in cc channel")
		respAllowedBalance, err := hlfProxy.Query("cc", "allowedBalanceOf", user.UserAddressBase58Check, "FIAT")
		sCtx.Assert().NoError(err)
		sCtx.Require().Equal("\"1\"", string(respAllowedBalance.Payload))
	})
	return swapBeginTxID, swapDoneTxID
}
