package config

import (
	"fmt"
	"go-rest-api/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Populate(db *gorm.DB) {
	if db.First(&models.Taxes{}).Error == nil {
		fmt.Println("Taxes table already populated, skipping...")
	} else {
		taxes := []*models.Taxes{
			{Year: 2021, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 3035.71428571},
			{Year: 2022, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 3063},
			{Year: 2023, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 3450},
		}

		res_tax := db.Create(&taxes)
		if res_tax.Error != nil {
			panic("Failed to create rows in taxes")
		}
	}

	if db.First(&models.PersonType{}).Error == nil {
		fmt.Println("PersonType table already populated, skipping...")
	} else {
		people := []*models.PersonType{
			{SocialStatus: "pensioner", IpnRate: 10, OpvAdd: false, VosmsAdd: false},
			{SocialStatus: "oppv", IpnRate: 8.8, OpvAdd: true, VosmsAdd: true},
			{SocialStatus: "mother", IpnRate: 9, OpvAdd: true, VosmsAdd: false},
			{SocialStatus: "disabled", IpnRate: 9, OpvAdd: true, VosmsAdd: false},
			{SocialStatus: "disabled_perm", IpnRate: 10, OpvAdd: false, VosmsAdd: false},
			{SocialStatus: "student", IpnRate: 9, OpvAdd: true, VosmsAdd: false},
			{SocialStatus: "asthub", IpnRate: 0, OpvAdd: true, VosmsAdd: true},
			{SocialStatus: "nonres_staff", IpnRate: 10, OpvAdd: false, VosmsAdd: false},
			{SocialStatus: "nonres_contr", IpnRate: 20, OpvAdd: false, VosmsAdd: false},
			{SocialStatus: "resident", IpnRate: 8.8, OpvAdd: true, VosmsAdd: true},
		}
		res_people := db.Create(&people)
		if res_people.Error != nil {
			panic("Failed to create rows in people")
		}
	}
}

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	if err != nil {
		panic("Failed to connect to database")
	}
	database.AutoMigrate(&models.Taxes{})
	database.AutoMigrate(&models.PersonType{})

	return database
}
