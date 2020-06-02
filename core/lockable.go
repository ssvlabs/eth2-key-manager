package core

type OperationType uint8
const (
	Creation 	OperationType = 0
	Sign	 	OperationType = 1
	FetchData 	OperationType = 2
)

// implements lock and unlock policies
type LockablePolicy interface {
	// Depending on the operation type, will return if should lock
	LockAfterOperation(op OperationType) bool
}