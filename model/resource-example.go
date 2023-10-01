package model

import "gorm.io/gorm"

type Resource struct {
	UserID int `gorm:"type:VARCHAR(128) NOT NULL;comment:用户主键" json:"userId"`

	Name string `gorm:"type:VARCHAR(128) NOT NULL;comment:名称" json:"name"`
	URL  string `gorm:"type:VARCHAR(128) NOT NULL;comment:资源URL" json:"url"`

	BaseModel
}

func (Resource) TableName() string {
	return "resource"
}

func (e *Resource) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func (e *Resource) BeforeUpdate(_ *gorm.DB) error {
	return nil
}

func (e *Resource) AfterFind(_ *gorm.DB) error {
	return nil
}
