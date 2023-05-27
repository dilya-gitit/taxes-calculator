package main

import (
	"errors"
	"go-rest-api/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Deduction string

const (
	Deduction14  Deduction = "14"
	Deduction882 Deduction = "882"
)

type TaxCalculator struct {
	Salary        float64 `json:"salary,string"`
	Deduction     Deduction
	Year          int      `json:"year,string"`
	IsStaffMember bool     `json:"isStaffMember,string" validate:"required,bool"`
	IsResident    bool     `json:"isResident,string" validate:"required,bool"`
	SocialStatus  []string `json:"socialStatus,omitempty"`
}

func main() {
	router := gin.Default()
	config.ConnectDatabase()
	config.Populate()
	router.POST("/post", calculateTaxesHandler)
	router.Run(":8080")
}

func validateTaxCalculator(tc TaxCalculator) error {
	if tc.Salary <= 0 {
		return errors.New("salary is required and сan't be negative")
	}

	switch tc.Deduction {
	case Deduction14, Deduction882:
		// Valid deduction value
	default:
		return errors.New("deduction value can be either 14 or 882")
	}

	if tc.Year < 2021 || tc.Year > 2023 || tc.Year == 0 {
		return errors.New("year is required and should be in the range of 2021 to 2023")
	}

	return nil
}

func calculateTaxesHandler(c *gin.Context) {

	var input TaxCalculator
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if err := validateTaxCalculator(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate taxes based on input
	taxes := calculateTaxes(input)

	c.JSON(http.StatusOK, gin.H{"taxes": taxes})
}

func calculateTaxes(input TaxCalculator) float64 {
	// Perform tax calculation based on the provided data
	// You can implement your own logic here based on the tax rules in your jurisdiction
	// For simplicity, let's assume a basic tax calculation based on salary and year

	var taxRate float64

	if input.Year < 2022 {
		taxRate = 0.2
	} else {
		taxRate = 0.18
	}

	taxes := input.Salary * taxRate

	return input.Salary - taxes
}
