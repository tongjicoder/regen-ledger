package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/regen-network/regen-ledger/x/ecocredit/v3"
	"github.com/regen-network/regen-ledger/x/ecocredit/v3/base"
	types "github.com/regen-network/regen-ledger/x/ecocredit/v3/base/types/v1"
	"github.com/regen-network/regen-ledger/x/ecocredit/v3/simulation/utils"
)

const OpWeightMsgUpdateProjectMetadata = "op_weight_msg_update_project_metadata" //nolint:gosec

var TypeMsgUpdateProjectMetadata = sdk.MsgTypeURL(&types.MsgUpdateProjectMetadata{})

const WeightUpdateProjectMetadata = 30

// SimulateMsgUpdateProjectMetadata generates a MsgUpdateProjectMetadata with random values.
func SimulateMsgUpdateProjectMetadata(ak ecocredit.AccountKeeper, bk ecocredit.BankKeeper, qryClient types.QueryServer) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, sdkCtx sdk.Context, accs []simtypes.Account, _ string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		class, op, err := utils.GetRandomClass(sdkCtx, r, qryClient, TypeMsgUpdateProjectMetadata)
		if err != nil {
			return op, nil, err
		}

		ctx := sdk.WrapSDKContext(sdkCtx)
		project, op, err := getRandomProjectFromClass(ctx, r, qryClient, TypeMsgUpdateProjectMetadata, class.Id)
		if project == nil {
			return op, nil, err
		}

		admin, err := sdk.AccAddressFromBech32(project.Admin)
		if err != nil {
			return simtypes.NoOpMsg(ecocredit.ModuleName, TypeMsgUpdateProjectMetadata, err.Error()), nil, err
		}

		msg := &types.MsgUpdateProjectMetadata{
			Admin:       admin.String(),
			ProjectId:   project.Id,
			NewMetadata: simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 10, base.MaxMetadataLength)),
		}

		spendable, account, op, err := utils.GetAccountAndSpendableCoins(sdkCtx, bk, accs, admin.String(), TypeMsgUpdateProjectMetadata)
		if spendable == nil {
			return op, nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         sdkCtx,
			SimAccount:      *account,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      ecocredit.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return utils.GenAndDeliverTxWithRandFees(r, txCtx)
	}
}
