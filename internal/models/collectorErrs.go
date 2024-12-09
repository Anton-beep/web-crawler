package models

type CounterIsNegative struct {
}

// Error CounterIsNegative is an error that is returned
// when the collector counter is negative
func (CounterIsNegative) Error() string {
	return "can't assign this message because collector counter is negative"
}

var CollectorCounterIsNegative CounterIsNegative
