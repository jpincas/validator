package validator

import (
	"testing"
	"time"
)

type maybeValid struct {
	validities []bool
}

func (m maybeValid) GetValidations() (vs []Validation) {
	for _, validity := range m.validities {
		x := validity //Important to cache this
		vs = append(vs, func(v Validateable) ValidationResult {
			//Since there's no real logic going on here, add a 1ms pause to simulate
			//a validation that actually does something
			time.Sleep(1 * time.Millisecond)
			return ValidationResult{x, "failed"}
		})
	}
	return
}

func TestOnePass(t *testing.T) {
	m := maybeValid{validities: []bool{true}}
	if validates := ValidateSerial(m); !validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); !validates.Pass {
		t.Fail()
	}
}

func TestMultiplePasses(t *testing.T) {
	m := maybeValid{validities: []bool{true, true, true, true, true}}
	if validates := ValidateSerial(m); !validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); !validates.Pass {
		t.Fail()
	}
}

func TestOneFail(t *testing.T) {
	m := maybeValid{validities: []bool{false}}
	if validates := ValidateSerial(m); validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); validates.Pass {
		t.Fail()
	}
}

func TestMultipleFails(t *testing.T) {
	m := maybeValid{validities: []bool{false, false, false, false}}
	if validates := ValidateSerial(m); validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); validates.Pass {
		t.Fail()
	}
}

func TestMixedPassFailFirst(t *testing.T) {
	m := maybeValid{validities: []bool{false, true, true}}
	if validates := ValidateSerial(m); validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); validates.Pass {
		t.Fail()
	}
}

func TestMixedPassFailMiddle(t *testing.T) {
	m := maybeValid{validities: []bool{true, false, true}}
	if validates := ValidateSerial(m); validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); validates.Pass {
		t.Fail()
	}
}

func TestMixedPassFailLast(t *testing.T) {
	m := maybeValid{validities: []bool{true, true, false}}
	if validates := ValidateSerial(m); validates.Pass {
		t.Fail()
	}
	if validates := ValidateParallel(m); validates.Pass {
		t.Fail()
	}
}

//Benchmarks

var vs = append((append(makeBoolSlice(100, true), false)), makeBoolSlice(100, true)...)

func BenchmarkMultipleSerial(b *testing.B) {

	m := maybeValid{validities: vs}
	for i := 0; i < b.N; i++ {
		ValidateSerial(m)
	}
}

func BenchmarkMultipleParallel(b *testing.B) {
	m := maybeValid{validities: vs}
	for i := 0; i < b.N; i++ {
		ValidateParallel(m)
	}
}

//Examples

//Utilities

func makeBoolSlice(n int, b bool) (s []bool) {
	for i := 1; i <= n; i++ {
		s = append(s, b)
	}
	return
}
