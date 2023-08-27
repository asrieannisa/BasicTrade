package models

import (
	"BasicTrade/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Admins struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UUID       string     `gorm:"not null" json:"uuid"`
	Name       string     `gorm:"not null" json:"name" form:"name" valid:"required~Your name is required"`
	Email      string     `gorm:"not null" json:"email" form:"email" valid:"required~Your email is required, email~Invalid email format"`
	Password   string     `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Products   []Products `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignkey:Admin_ID" json:"products"`
	Created_at *time.Time `json:"created_at,omitempty"`
	Updated_at *time.Time `json:"updated_at,omitempty"`
}

func (u *Admins) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPass(u.Password)

	err = nil
	return
}

func (b *Admins) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
