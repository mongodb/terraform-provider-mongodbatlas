package retrystrategy

const (
	RetryStrategyPendingState    = "PENDING"
	RetryStrategyCompletedState  = "COMPLETED"
	RetryStrategyErrorState      = "ERROR"
	RetryStrategyPausedState     = "PAUSED"
	RetryStrategyUpdatingState   = "UPDATING"
	RetryStrategyInitiatingState = "INITIATING"
	RetryStrategyIdleState       = "IDLE"
	RetryStrategyActiveState     = "ACTIVE"
	RetryStrategyDeletedState    = "DELETED"

	RetryStrategyPendingAcceptanceState = "PENDING_ACCEPTANCE"
	RetryStrategyPendingRecreationState = "PENDING_RECREATION"
)
