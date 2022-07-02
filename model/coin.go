package model

import (
	"gorm.io/gorm"
)

type Coin struct {
	*gorm.Model
	ChainID  uint64
	Name     string `gorm:"type:char(32)"`
	Decimal  uint
	Contract string `gorm:"type:char(256)"`
	GasPrice string `gorm:"type:char(64)"`
	Gaslimit string `gorm:"type:char(64)"`
	BaseFee  string `gorm:"type:char(64)"`
	Fee      string `gorm:"type:char(64)"`
}
