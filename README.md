[![GoDoc](https://godoc.org/github.com/jpincas/validator?status.svg)](https://godoc.org/github.com/jpincas/validator)

# Simple Validation Runner

A simple package to help run logical validations on structs.  Just define your validations as functions returning a Boolean result and failure message, and specify a `GetValidations()` method on your struct listing which validations to run.  You then use `ValidateSerial()` or `ValidateParallel()` to execute your validations.

## Example

```golang
package validator

import (
	"errors"
	"fmt"
)

//First, we define our customer struct.
type customer struct {
	age            int
	agreedToTerms  bool
	accountBlocked bool
	email          string
}

//Next define some simple logical validations with accompanying error messages.
//You could have one or more validations per validation function.
//Here we validate only age...
func isOver18(v Validateable) error {
	if v.(customer).age < 18 {
		return errors.New("Customer is under age")
	}
	return nil
}

//..whereas here we combine two validations in one function
//with
func isActive(v Validateable) error {
	if !v.(customer).agreedToTerms || v.(customer).accountBlocked {
		return errors.New("Customer account inactive")
	}
	return nil
}

func hasEmail(v Validateable) error {
	if v.(customer).email == "" {
		return errors.New("Customer has no email")
	}
	return nil
}

//Finally, we list the set of validations to perform for 'customer',
//thereby fulfilling the 'Validateable' interface
func (c customer) GetValidations() []Validation {
	return []Validation{isOver18, isActive, hasEmail}
}

//Now we can create some customers and validate them.
func Example() {
	validCustomer := customer{37, true, false, "email@email.com"}
	underAgeCustomer := customer{16, true, false, "email@email.com"}
	disagreeingCustomer := customer{21, false, false, "email@email.com"}
	blockedCustomer := customer{21, true, true, "email@email.com"}
	noEmailCustomer := customer{21, true, false, ""}
	//We can choose to run validations in serial or parallel.
	//As a rule of thumb, parallel will be faster for any validations that involve
	//complex computations or data access.  The included benchmarks shows parallel mode
	//to be roughly 100X faster for a simulated 1ms validation delay
	r1 := ValidateSerial(validCustomer)
	r2 := ValidateParallel(underAgeCustomer)
	r3 := ValidateSerial(disagreeingCustomer)
	r4 := ValidateParallel(blockedCustomer)
	r5 := ValidateSerial(noEmailCustomer)

	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println(r3)
	fmt.Println(r4)
	fmt.Println(r5)

	// Output:
	// <nil>
	// Customer is under age
	// Customer account inactive
	// Customer account inactive
	// Customer has no email
}

```