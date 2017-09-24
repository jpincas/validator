package validator

import "fmt"

// Validateable interface
type Validateable interface {
	GetValidations() []Validation
}

type Validation func(v Validateable) ValidationResult

type ValidationResult struct {
	Pass           bool
	FailureMessage string
}

func (vr ValidationResult) String() string {
	if vr.Pass {
		return "validation passed"
	}
	return fmt.Sprintf("validation failed: %s", vr.FailureMessage)
}

func ValidateSerial(v Validateable) ValidationResult {
	validations := v.GetValidations()
	for _, validation := range validations {
		if result := validation(v); !result.Pass {
			return ValidationResult{false, result.FailureMessage}
		}
	}
	return ValidationResult{true, ""}
}

func ValidateParallel(v Validateable) ValidationResult {
	//Set up a buffered channel of size equal to the number of validations
	validations := v.GetValidations()
	noOfValidations := len(validations)
	results := make(chan ValidationResult, noOfValidations)
	//Run each validation in parallel
	for _, validation := range validations {
		go func(thisValidation Validation) {
			results <- thisValidation(v)
		}(validation)
	}
	//Read the results as they come in
	//Return false as soon as a failed validation is reported
	for i := 0; i < noOfValidations; i++ {
		r := <-results
		if !r.Pass {
			return ValidationResult{false, r.FailureMessage}
		}
	}
	//Otherwise all validations have passed
	return ValidationResult{true, ""}
}
