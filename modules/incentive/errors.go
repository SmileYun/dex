package incentive

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceIncentive sdk.CodespaceType = "incentive"

	// 701 ～ 799
	CodeInvalidIncentiveBlockInterval sdk.CodeType = 701
	CodeInvalidDefaultRewardPerBlock  sdk.CodeType = 702
	CodeInvalidPlanHeight             sdk.CodeType = 703
	CodeInvalidRewardPerBlock         sdk.CodeType = 704
	CodeInvalidTotalIncentive         sdk.CodeType = 705
)