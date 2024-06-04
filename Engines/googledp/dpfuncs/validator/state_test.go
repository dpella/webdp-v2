package validator

import (
	"googledp/entities"
	"testing"
)

func TestValidator1(t *testing.T) {
	testArr := []string{entities.FILTER, entities.BIN, entities.MEASUREMENT}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != true {
		t.Fatalf("test failed")
	}
}

func TestValidator2(t *testing.T) {
	testArr := []string{entities.FILTER, entities.MEASUREMENT}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != true {
		t.Fatalf("test failed")
	}
}

func TestValidator3(t *testing.T) {
	testArr := []string{entities.BIN, entities.MEASUREMENT}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != true {
		t.Fatalf("test failed")
	}
}

func TestValidator4(t *testing.T) {
	testArr := []string{entities.MEASUREMENT}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != true {
		t.Fatalf("test failed")
	}
}

func TestValidator5(t *testing.T) {
	testArr := []string{entities.FILTER, entities.BIN}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != false {
		t.Fatalf("test failed")
	}
}

func TestValidator6(t *testing.T) {
	testArr := []string{entities.FILTER}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != false {
		t.Fatalf("test failed")
	}
}

func TestValidator7(t *testing.T) {
	testArr := []string{entities.BIN}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != false {
		t.Fatalf("test failed")
	}
}

func TestValidator8(t *testing.T) {
	testArr := []string{entities.BIN, entities.FILTER, entities.MEASUREMENT}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != false {
		t.Fatalf("test failed")
	}
}

func TestValidator9(t *testing.T) {
	testArr := []string{entities.BIN, entities.FILTER}
	validator := NewSMValidator()
	if validator.VerifyInputs(testArr) != false {
		t.Fatalf("test failed")
	}
}
