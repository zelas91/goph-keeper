package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"regexp"
)

func NewValidator(log logger.Logger) *validator.Validate {
	valid := validator.New()
	if err := valid.RegisterValidation("expired_credit_card", expiredCreditCard); err != nil {
		log.Errorf("Register custom validator err:%v", err)
	}
	return valid
}
func expiredCreditCard(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	dateRegex := regexp.MustCompile(`^(\d{2})/(\d{2})$`)

	if !dateRegex.MatchString(date) {
		return false
	}

	month := dateRegex.ReplaceAllString(date, "$1")
	year := dateRegex.ReplaceAllString(date, "$2")

	if month > "12" || month == "00" {
		return false
	}

	if year < "19" || year > "99" {
		return false
	}

	return true
}
