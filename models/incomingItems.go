package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type IncomingItems struct {
	GormModel
	Qty        uint8      `gorm:"not null" json:"qty" form:"qty" valid:"required~Your quantity of incoming is required"`
	IncomingAt CustomTime `gorm:"not null" json:"incoming_at" form:"incoming_at" valid:"required~Your incoming at of incoming is required"`
	UserID     uint       `gorm:"not null" json:"user_id" form:"user_id" valid:"required~Your user id is required"`
	ProductID  uint       `gorm:"not null" json:"product_id" form:"product_id" valid:"required~Your product id is required"`
	Products   *Products  `gorm:"foreignKey:ProductID;references:ID" json:"products"`
	Users      *Users     `gorm:"foreignKey:UserID;references:ID" json:"users"`
}

func (p *IncomingItems) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (p *IncomingItems) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(p)

	if errUpdate != nil {
		err = errUpdate
		return
	}

	err = nil
	return
}
