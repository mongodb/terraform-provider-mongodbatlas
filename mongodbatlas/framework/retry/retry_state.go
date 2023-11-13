package retry

const (
	RetryStrategyPendingState   = "PENDING"
	RetryStrategyCompletedState = "COMPLETED"
	RetryStrategyErrorState     = "ERROR"
	RetryStrategyPausedState    = "PAUSED"
	RetryStrategyUpdatingState  = "UPDATING"
	RetryStrategyIdleState      = "IDLE"
	RetryStrategyDeletedState   = "DELETED"
)
