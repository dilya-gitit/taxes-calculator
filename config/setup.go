package config

import (
	"go-rest-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Populate() {
	taxes := []*models.Taxes{
		{Year: 2021, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 2917},
		{Year: 2022, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 3063},
		{Year: 2023, TaxTypeIpn: 10, TaxTypeOpv: 10, TaxTypeVosms: 2, Mrp: 3450},
	}

	res_tax := DB.Create(&taxes) // pass a slice to insert multiple row
	if res_tax.Error != nil {
		panic("Failed to create rows in taxes")
	}
	people := []*models.People{
		{SocialStatus: "Pensioner", IpnRate: 10, OpvExempt: false, VosmsExempt: false},
		{SocialStatus: "OPPV", IpnRate: 8.8, OpvExempt: true, VosmsExempt: true},
		{SocialStatus: "Mother", IpnRate: 9, OpvExempt: true, VosmsExempt: false},
		{SocialStatus: "Disabled", IpnRate: 9, OpvExempt: true, VosmsExempt: false},
		{SocialStatus: "DisabledInf", IpnRate: 10, OpvExempt: false, VosmsExempt: false},
		{SocialStatus: "Student", IpnRate: 9, OpvExempt: true, VosmsExempt: false},
		{SocialStatus: "Asthub_staff", IpnRate: 0, OpvExempt: true, VosmsExempt: true},
		{SocialStatus: "Asthub_contr", IpnRate: 8.8, OpvExempt: true, VosmsExempt: true},
	}
	res_people := DB.Create(&people)
	if res_people.Error != nil {
		panic("Failed to create rows in people")
	}

}
func ConnectDatabase() {
	/*dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})*/
	database, err := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5432/postgres"))
	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&models.Taxes{})
	database.AutoMigrate(&models.People{})

	DB = database
}
