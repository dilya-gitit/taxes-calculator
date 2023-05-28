package models

import "gorm.io/gorm"

// Taxes model
type Taxes struct {
	gorm.Model
	Year         int     `gorm:"primaryKey"`
	TaxTypeIpn   float64 `gorm:"column:tax_type_ipn"`
	TaxTypeOpv   float64 `gorm:"column:tax_type_opv"`
	TaxTypeVosms float64 `gorm:"column:tax_type_vosms"`
	Mrp          float64 `gorm:"column:mrp"`
}

// PersonType model
type PersonType struct {
	gorm.Model
	SocialStatus string  `gorm:"primaryKey"`
	IpnRate      float64 `gorm:"column:ipn_rate"`
	OpvAdd       bool    `gorm:"column:opv_add"`
	VosmsAdd     bool    `gorm:"column:vosms_add"`
}
