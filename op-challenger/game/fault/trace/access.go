package trace

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

type ProviderCreator func(ctx context.Context, pre types.Claim, post types.Claim) (types.TraceProvider, error)

func NewSimpleTraceAccessor(trace types.TraceProvider) *Accessor {
	selector := func(_ context.Context, _ types.Game, _ types.Claim, _ types.Position) (types.TraceProvider, error) {
		return trace, nil
	}
	return &Accessor{selector}
}

type Accessor struct {
	selector func(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error)
}

func (t *Accessor) Get(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (common.Hash, error) {
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return common.Hash{}, err
	}
	return provider.Get(ctx, pos)
}

func (t *Accessor) GetStepData(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	return provider.GetStepData(ctx, pos)
}

var _ types.TraceAccessor = (*Accessor)(nil)
