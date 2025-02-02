package bot

import (
	"context"
	"fmt"

	"github.com/MixinNetwork/mixin/common"
)

func CreateObjectStorageTransaction(ctx context.Context, recipients []*TransactionRecipient, extra []byte, traceId string, references []string, limit string, u *SafeUser) (*SequencerTransactionRequest, error) {
	if len(extra) > common.ExtraSizeStorageCapacity {
		return nil, fmt.Errorf("too large extra %d > %d", len(extra), common.ExtraSizeStorageCapacity)
	}
	step := common.NewIntegerFromString(common.ExtraStoragePriceStep)
	amount := step.Mul(len(extra)/common.ExtraSizeStorageStep + 1)
	if limit != "" {
		strl := common.NewIntegerFromString(limit)
		if strl.Cmp(amount) > 0 {
			amount = strl
		}
	}
	addr := common.NewAddressFromSeed(make([]byte, 64))
	mix := NewMainnetMixAddress([]string{addr.String()}, 1)
	mix.Threshold = 64
	rec := []*TransactionRecipient{{
		MixAddress: mix,
		Amount:     amount.String(),
	}}
	if len(recipients) > 0 {
		rec = append(rec, recipients...)
	}
	return SendTransaction(ctx, common.XINAssetId.String(), rec, traceId, extra, references, u)
}
