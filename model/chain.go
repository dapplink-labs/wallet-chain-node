package model

import (
	"gorm.io/gorm"
)

type Chain struct {
	*gorm.Model
	Name        string `gorm:"type:char(32)"`
	Index       uint
	RpcUrl      string `gorm:"type:char(256)"`
	RpcUser     string `gorm:"type:char(64)"`
	RpcPassword string `gorm:"type:char(64)"`
	ScanUrl     string `gorm:"type:char(256)"`
	ScanApiKey  string `gorm:"type:char(256)"`
}
