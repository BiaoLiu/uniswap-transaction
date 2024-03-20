package data

import "time"

type TransactionType string

const (
	Mint TransactionType = "mint"
	Burn TransactionType = "burn"
	Swap TransactionType = "swap"
)

type Transaction struct {
	ID         int64           `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Hash       string          `gorm:"column:hash;not null" json:"hash"`
	Token0     string          `gorm:"column:token0;not null" json:"token0"`
	Token1     string          `gorm:"column:token1;not null" json:"token1"`
	Type       TransactionType `gorm:"column:type;not null" json:"type"`
	Amount0    int64           `gorm:"column:amount0;not null" json:"amount0"`
	Amount1    int64           `gorm:"column:amount1;not null" json:"amount1"`
	Amount0In  int64           `gorm:"column:amount0_in;not null" json:"amount0_in"`
	Amount1Out int64           `gorm:"column:amount1_out;not null" json:"amount1_out"`
	From       string          `gorm:"column:from;not null" json:"from"`
	To         string          `gorm:"column:to;not null" json:"to"`
	CreatedAt  time.Time       `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (*Transaction) TableName() string {
	return "transaction"
}
