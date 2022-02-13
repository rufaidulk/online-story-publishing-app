package models

import "time"

const PaymentCompleted int8 = 10
const PaymentFailed int8 = 11

type PaymentTransaction struct {
	Id          int64
	RefNo       string
	Amount      float64
	Type        string
	TxnCategory string
	Status      int8
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPaymentTransaction() *PaymentTransaction {
	return &PaymentTransaction{}
}

func (p *PaymentTransaction) SetPaymentCompleted() {
	p.Status = PaymentCompleted
}

func (p *PaymentTransaction) SetPaymentFailed() {
	p.Status = PaymentFailed
}
