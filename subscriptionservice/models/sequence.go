package models

type Sequence struct {
	Id    int64
	Type  string
	SeqNo int64
}

const SubscriptionTxn string = "subscriptionTxn"

func NewSequence() *Sequence {
	return &Sequence{}
}
