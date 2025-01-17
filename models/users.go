package models

import (
	"inventoryapp/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Users struct {
	GormModel
	Username string `gorm:"unique;not null;uniqueIndex" json:"username" form:"username" valid:"required~Your username is required"`
	Email    string `gorm:"unique;not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Your password must be at least 6 characters"`
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPass(u.Password)
	err = nil
	return
}
