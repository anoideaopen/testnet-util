package transfer

import (
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	utils "github.com/atomyze-foundation/atomyze-util"
)

func ChannelTransferByCustomer(t provider.T, hlfProxy *utils.HlfProxyService, channelFrom string, user utils.User, transferArgs []string) {
	t.WithNewStep("Signing transfer args and invoke channelTransferByCustomer then checking balance of head channel", func(sCtx provider.StepCtx) {
		sa, err := utils.Sign(user.UserEd25519PrivateKey, user.UserEd25519PublicKey, channelFrom, channelFrom, "channelTransferByCustomer", transferArgs)
		sCtx.Require().NoError(err)

		_, err = hlfProxy.Invoke(channelFrom, "channelTransferByCustomer", sa...)
		sCtx.Require().NoError(err)

		time.Sleep(utils.BatchTransactionTimeout)
	})
}

func ChannelTransferByAdmin(t provider.T, hlfProxy *utils.HlfProxyService, issuer utils.Issuer, channelFrom string, transferArgs []string) {
	t.WithNewStep("Signing transfer args and invoke channelTransferByAdmin then checking balance of head channel", func(sCtx provider.StepCtx) {
		sa, err := utils.Sign(issuer.IssuerEd25519PrivateKey, issuer.IssuerEd25519PublicKey, channelFrom, channelFrom, "channelTransferByAdmin", transferArgs)
		sCtx.Require().NoError(err)

		_, err = hlfProxy.Invoke(channelFrom, "channelTransferByAdmin", sa...)
		sCtx.Require().NoError(err)

		time.Sleep(utils.BatchTransactionTimeout)
	})
}

func ChannelTransferFrom(t provider.T, hlfProxy *utils.HlfProxyService, channel string, transferID string) string {
	var form string
	t.WithNewStep("Getting a transfer record from outgoing channel with channelTransferFrom", func(sCtx provider.StepCtx) {
		resp, err := hlfProxy.Invoke(channel, "channelTransferFrom", transferID)
		t.Require().NoError(err)
		form = string(resp.Payload)
	})
	return form
}

func CreateCCTransferTo(t provider.T, hlfProxy *utils.HlfProxyService, channelTo string, form string) {
	t.WithNewStep("create cc transfer", func(sCtx provider.StepCtx) {
		transferArgs := []string{form}
		_, err := hlfProxy.Invoke(channelTo, "createCCTransferTo", transferArgs...)
		t.Require().NoError(err)
		time.Sleep(utils.BatchTransactionTimeout)
	})
}

func ChannelTransferTo(t provider.T, hlfProxy *utils.HlfProxyService, channelTo string, transferID string) {
	t.WithNewStep("channel transfer", func(sCtx provider.StepCtx) {
		a := []string{transferID}
		_, err := hlfProxy.Invoke(channelTo, "channelTransferTo", a...)
		t.Require().NoError(err)
	})
}

func CheckAllowedBalanceEqual(t provider.T, hlfProxy *utils.HlfProxyService, userAddressBase58Check string, channel string, token string, amount string) {
	t.WithNewStep("Checking that balance equal "+amount, func(sCtx provider.StepCtx) {
		respGetBalance, err := hlfProxy.Query(channel, "allowedBalanceOf", userAddressBase58Check, token)
		sCtx.Require().NoError(err)
		sCtx.Require().Equal("\""+amount+"\"", string(respGetBalance.Payload))
	})
}

func CommitCCTransferFrom(t provider.T, hlfProxy *utils.HlfProxyService, channelFrom string, transferID string) {
	t.WithNewStep("commit CC transfer from", func(sCtx provider.StepCtx) {
		a := []string{transferID}
		_, err := hlfProxy.Invoke(channelFrom, "commitCCTransferFrom", a...)
		t.Require().NoError(err)
	})
}

func DeleteCCTransferTo(t provider.T, hlfProxy *utils.HlfProxyService, channelTo string, transferID string) {
	t.WithNewStep("dalete CC transfer to", func(sCtx provider.StepCtx) {
		a := []string{transferID}
		_, err := hlfProxy.Invoke(channelTo, "deleteCCTransferTo", a...)
		t.Require().NoError(err)
	})
}

func DeleteCCTransferFrom(t provider.T, hlfProxy *utils.HlfProxyService, channelFrom string, transferID string) {
	t.WithNewStep("delete CC transfer from", func(sCtx provider.StepCtx) {
		a := []string{transferID}
		_, err := hlfProxy.Invoke(channelFrom, "deleteCCTransferFrom", a...)
		t.Require().NoError(err)
	})
}

func CancelCCTransferFrom(t provider.T, hlfProxy *utils.HlfProxyService, channelFrom string, transferID string) {
	t.WithNewStep("cancel CC transfer from", func(sCtx provider.StepCtx) {
		a := []string{transferID}
		_, err := hlfProxy.Invoke(channelFrom, "cancelCCTransferFrom", a...)
		t.Require().NoError(err)
		time.Sleep(utils.BatchTransactionTimeout)
	})
}

func ChannelTransfersFrom(t provider.T, hlfProxy *utils.HlfProxyService, channelFrom string, pageSize string, bookmark string) []byte {
	var payload []byte
	t.WithNewStep("channel transfer from", func(sCtx provider.StepCtx) {
		a := []string{pageSize, bookmark}
		resp, err := hlfProxy.Invoke(channelFrom, "channelTransfersFrom", a...)
		t.Require().NoError(err)
		payload = resp.Payload
	})
	return payload
}
