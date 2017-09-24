// Package validator provides a simple way to help run logical validations on structs.
// Just define your validations as functions returning a Boolean result and failure message,
// and specify a `GetValidations()` method on your struct listing which validations to run.
// You then use `ValidateSerial()` or `ValidateParallel()` to execute your validations.
package validator

import "fmt"

// Validateable interface is implemented by supplying a list of validation functions to be executed
type Validateable interface {
	GetValidations() []Validation
}

//Validation is a validation function to be executed as part of the validation
type Validation func(v Validateable) ValidationResult

//ValidationResult is simply true/false combined with an error message to be used in case of failure
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

//ValidateSerial executes validations functions one-at-a-time in order
//It will exit on first fail, returning the corresponding failure message
func ValidateSerial(v Validateable) ValidationResult {
	validations := v.GetValidations()
	for _, validation := range validations {
		if result := validation(v); !result.Pass {
			return ValidationResult{false, result.FailureMessage}
		}
	}
	return ValidationResult{true, ""}
}

//ValidateParallel executes validations in parallel.
//It reads the validatio results as they become available,
//and will exit immediately on receiving a fail
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
