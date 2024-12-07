package models

type CounterIsNegative struct {
}

func (CounterIsNegative) Error() string {
	return "can't assign this message because collector counter is negative"
}

var CollectorCounterIsNegative CounterIsNegative
