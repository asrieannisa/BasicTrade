package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Variants struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UUID         string `gorm:"not null" json:"uuid"`
	Variant_name string `gorm:"not null" json:"variant_name" form:"variant_name" valid:"required~Variant name is required"`
	Quantity     int    `gorm:"not null" json:"quantity"`
	Product_ID   uint   `gorm:"column:product_id"`
	Products     *Products
	//Products     *Products `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL; foreignKey:Product_ID" json:"products"`
	Created_at time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	Updated_at time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

func (b *Variants) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (b *Variants) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(b)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
