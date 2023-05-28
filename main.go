package main

import (
	"errors"
	"go-rest-api/config"
	"go-rest-api/models"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Deduction string

const (
	Deduction14  Deduction = "14"
	Deduction882 Deduction = "882"
)

type Request struct {
	Salary        float64 `json:"salary,string"`
	Deduction     Deduction
	Year          int      `json:"year,string" validate:"required,int"`
	IsStaffMember bool     `json:"is_staff_member,string" validate:"required,bool"`
	IsResident    bool     `json:"is_resident,string" validate:"required,bool"`
	SocialStatus  []string `json:"social_statuses,omitempty"`
}

func main() {
	router := gin.Default()
	db := config.ConnectDatabase()
	config.Populate(db)
	router.POST("/calculate-tax", calculateTaxesHandler(db))
	router.Run(":8080")
}

func validateRequest(tc Request) error {
	if tc.Salary <= 0 {
		return errors.New("salary is required and can't be negative")
	}

	if tc.Year < 2021 || tc.Year > 2023 || tc.Year == 0 {
		return errors.New("year is required and should be in the range of 2021 to 2023 inclusive")
	}

	if !tc.IsResident && len(tc.SocialStatus) > 0 {
		return errors.New("a foreigner can't be assigned a social status")
	}
	if (!tc.IsStaffMember || !tc.IsResident) && tc.Deduction == Deduction14 {
		return errors.New("a contractor or a foreigner can't get a 14 MRP Deduction")
	}
	return nil
}

func calculateTaxesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input Request
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validateRequest(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Calculate taxes based on input
		taxSummary, deduction := calculateTaxes(input, db)
		netSalary, resTaxSummary := calculateSalary(taxSummary, input.Salary, deduction)

		response := gin.H{
			"net_salary": int(netSalary),
		}

		if resTaxSummary.IPN != 0 {
			response["IPN"] = int(resTaxSummary.IPN)
		}
		if resTaxSummary.OPV != 0 {
			response["OPV"] = int(resTaxSummary.OPV)
		}
		if resTaxSummary.VOSMS != 0 {
			response["VOSMS"] = int(resTaxSummary.VOSMS)
		}
		c.JSON(http.StatusOK, response)
	}
}

type TaxSummary struct {
	IPN   float64
	OPV   float64
	VOSMS float64
}

func calculateTaxes(input Request, db *gorm.DB) (TaxSummary, float64) {
	var people []models.PersonType
	var taxes models.Taxes
	db.Where("year = ?", input.Year).Find(&taxes)
	taxSummary := TaxSummary{}
	deduction := 0.0

	switch input.Deduction {
	case Deduction14:
		deduction = 14 * taxes.Mrp
	case Deduction882:
		if input.IsResident {
			deduction = 882 * taxes.Mrp
		}
	}
	if !input.IsResident {
		if input.IsStaffMember {
			db.Where("social_status = ?", "nonres_staff").Find(&people)
		} else {
			db.Where("social_status = ?", "nonres_contr").Find(&people)
		}
		taxSummary = calculateTaxForStatusType(people[0], taxes)
	} else {
		if len(input.SocialStatus) == 0 {
			db.Where("social_status = ?", "resident").Find(&people)
			taxSummary = calculateTaxForStatusType(people[0], taxes)
		} else {
			db.Where("social_status IN ?", input.SocialStatus).Find(&people)
			taxSummary.IPN = 10000

			// because some taxes override another
			priorityOrder := []string{
				"pensioner",
				"disabled_perm",
				"disabled",
				"mother",
				"student",
				"asthub",
				"oppv"}
			for _, status := range priorityOrder {
				for _, person := range people {
					if person.SocialStatus == status {
						taxSummary = calculateTaxForStatusType(person, taxes)
						return taxSummary, deduction
					}
				}
			}
		}
	}
	return taxSummary, deduction
}

func calculateTaxForStatusType(person models.PersonType, taxes models.Taxes) TaxSummary {
	ipn := person.IpnRate
	opv := 0.0
	vosms := 0.0

	if person.OpvAdd {
		opv = taxes.TaxTypeOpv
	}
	if person.VosmsAdd {
		vosms = taxes.TaxTypeVosms
	}
	return TaxSummary{ipn, opv, vosms}
}

func calculateSalary(taxSummary TaxSummary, salary float64, deduction float64) (float64, TaxSummary) {
	newTaxes := TaxSummary{
		IPN:   salary * taxSummary.IPN * 0.01,
		OPV:   salary * taxSummary.OPV * 0.01,
		VOSMS: salary * taxSummary.VOSMS * 0.01,
	}

	if deduction != 0 && newTaxes.IPN != 0 {
		taxesTotal := newTaxes.OPV + newTaxes.VOSMS
		newTaxes.IPN = math.Max(((salary - taxesTotal - deduction) * 0.1), 0)
		return salary - newTaxes.IPN - taxesTotal, newTaxes
	}

	netSalary := salary - newTaxes.IPN - newTaxes.OPV - newTaxes.VOSMS
	return netSalary, newTaxes
}
