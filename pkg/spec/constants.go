package spec

/*
Phase0
*/

const (
	MaxEffectiveInc        = 32
	BaseRewardFactor       = 25
	BaseRewardPerEpoch     = 4
	EffectiveBalanceInc    = 1000000000
	SlotsPerEpoch          = 16
	ProposerRewardQuotient = 8
	SlotsPerHistoricalRoot = 8192
	SlotSeconds            = 5
	EpochSlots             = 16
)

/*
Altair
*/
const (
	// spec weight constants
	TimelySourceWeight = 14
	TimelyTargetWeight = 26
	TimelyHeadWeight   = 14

	SyncRewardWeight  = 2
	ProposerWeight    = 8
	WeightDenominator = 64
	SyncCommitteeSize = 512
)

var (
	ParticipatingFlagsWeight = [3]int{TimelySourceWeight, TimelyTargetWeight, TimelyHeadWeight}
)

type ModelType int8

const (
	BlockModel ModelType = iota
	BlockDropModel
	OrphanModel
	EpochModel
	EpochDropModel
	PoolSummaryModel
	ProposerDutyModel
	ProposerDutyDropModel
	ValidatorLastStatusModel
	ValidatorRewardsModel
	ValidatorRewardDropModel
	WithdrawalModel
	WithdrawalDropModel
	TransactionsModel
	TransactionDropModel
	ReorgModel
	FinalizedCheckpointModel
)

type ValidatorStatus int8

const (
	QUEUE_STATUS ValidatorStatus = iota
	ACTIVE_STATUS
	EXIT_STATUS
	SLASHED_STATUS
)
