package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Products struct {
	GormModel
	Name  string `gorm:"not null" json:"name" form:"name" valid:"required~Your product name is required"`
	Stock uint8  `json:"stock" form:"stock"`
}

func (p *Products) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (p *Products) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(p)

	if errUpdate != nil {
		err = errUpdate
		return
	}

	err = nil
	return

}
