package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Products struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	UUID     string     `gorm:"not null" json:"uuid"`
	Name     string     `gorm:"not null" json:"name" form:"name" valid:"required~Product name is required"`
	ImageURL string     `json:"image_url"`
	Variants []Variants `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignkey:Product_ID" json:"variants"`
	Admin_ID uint       `gorm:"column:admin_id"`
	Admin    *Admins
	//Admin    *Admins    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:Admin_ID"`
	// Admins     *Admins    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:Admin_ID" json:"admins"`
	Created_at time.Time `gorm:"not null" json:"created_at"`
	Updated_at time.Time `gorm:"not null" json:"updated_at"`
}

func (b *Products) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (b *Products) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (b *Products) BeforeDelete(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
