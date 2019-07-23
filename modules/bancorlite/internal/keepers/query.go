package keepers

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryBancorInfo = "bancor-info"
)

// creates a querier for asset REST endpoints
func NewQuerier(keeper Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abcitypes.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBancorInfo:
			return queryBancorInfo(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("query symbol : " + path[0])
		}
	}
}

type QueryBancorInfoParam struct {
	Token string `json:"token"`
}

func queryBancorInfo(ctx sdk.Context, req abcitypes.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var param QueryBancorInfoParam
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &param); err != nil {
		return nil, sdk.NewError(types.CodeSpaceBancorlite, types.CodeUnMarshalFailed, "failed to parse param")
	}
	bi := keeper.Bik.Load(ctx, param.Token)

	biD := NewBancorInfoDisplay(bi)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, biD)
	if err != nil {
		return nil, sdk.NewError(types.CodeSpaceBancorlite, types.CodeMarshalFailed, "could not marshal result to JSON")
	}
	return bz, nil
}