package bankx

import (
	"github.com/coinexchain/dex/x/authx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}

}
func subActivatedFee(amt sdk.Coins, sub int64) (sdk.Coins, bool) {

	var found bool
	for i, coin := range amt {
		if coin.Denom != "cet" {
			continue
		}

		amt[i] = sdk.NewCoin("cet", coin.Amount.Sub(sdk.NewInt(sub)))
		found = true
	}
	return amt, found
}

func handleMsgSend(ctx sdk.Context, k Keeper, msg bank.MsgSend) sdk.Result {

	if !k.bk.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.bk.Codespace()).Result()
	}

	_, ok := k.axk.GetAccountX(ctx, msg.ToAddress)
	var amt sdk.Coins
	var found bool

	activatedFee := k.GetParam(ctx).ActivatedFee
	//toaccount doesn't exist yet
	if !ok {

		//check whether the first transfer contains cet

		amt, found = subActivatedFee(msg.Amount, activatedFee)
		if !found {
			return ErrorFirstTransferNotCET(CodeSpaceBankx).Result()
		}
		if !amt.IsValid() {
			return sdk.ErrInvalidCoins(amt.String()).Result()
		}

		//collect account activation fees
		//TODO: update token symbol rules for cetchain to allow "cet"
		k.fck.AddCollectedFees(ctx, sdk.NewCoins(sdk.NewCoin("cet,", sdk.NewInt(activatedFee))))
	}

	// sub the activatedfees from fromaddress
	fromAccount := k.ak.GetAccount(ctx, msg.FromAddress)
	oldCoins := fromAccount.GetCoins()
	newCoins, _ := subActivatedFee(oldCoins, activatedFee)
	fromAccount.SetCoins(newCoins)
	k.ak.SetAccount(ctx, fromAccount)

	//handle coins transfer
	t, err := k.bk.SendCoins(ctx, msg.FromAddress, msg.ToAddress, amt)

	if err != nil {
		return err.Result()
	}

	// new accountx for toaddress if needed
	if !ok {
		newAccountX := authx.NewAccountXWithAddress(msg.ToAddress)
		newAccountX.Activated = true
		k.axk.SetAccountX(ctx, newAccountX)
	}

	return sdk.Result{
		Tags: t,
	}
}
