package retry

const (
	RetryStrategyPendingState   = "PENDING"
	RetryStrategyCompletedState = "COMPLETED"
	RetryStrategyErrorState     = "ERROR"
	RetryStrategyPausedState    = "PAUSED"
	RetryStrategyPlanningState  = "PLANNING"
	RetryStrategyWorkingState   = "WORKING"
	RetryStrategyDeletedState   = "DELETED"
	RetryStrategyIdleState      = "IDLE"
)
