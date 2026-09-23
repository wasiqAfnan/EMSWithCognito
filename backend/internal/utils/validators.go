package utils

import (
	"fmt"
	"regexp"
	"strings"

	"awsems/internal/model"
)

// emailRegex is a simple RFC-5322-compatible email validator pattern.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateString trims whitespace and ensures the value is non-empty.
func ValidateString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	return nil
}

// ValidateEmpID checks that the employee ID:
//   - is non-empty
//   - starts with "EMP"
//   - has only digits after the prefix
//   - is at least 6 characters long (e.g. EMP001)
func ValidateEmpID(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return fmt.Errorf("empId must not be empty")
	}
	if !strings.HasPrefix(value, "EMP") {
		return fmt.Errorf("empId must start with 'EMP'")
	}
	suffix := value[3:]
	if len(suffix) == 0 {
		return fmt.Errorf("empId must contain digits after 'EMP'")
	}
	for _, ch := range suffix {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("empId must contain only digits after 'EMP'")
		}
	}
	if len(value) < 6 {
		return fmt.Errorf("empId must be at least 6 characters long")
	}
	return nil
}

// ValidateContactNo checks that the contact number:
//   - is non-empty
//   - contains only digits
//   - is between 10 and 12 digits long
func ValidateContactNo(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return fmt.Errorf("contactNo must not be empty")
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("contactNo must contain only digits")
		}
	}
	if len(value) < 10 {
		return fmt.Errorf("contactNo must be at least 10 digits")
	}
	if len(value) > 12 {
		return fmt.Errorf("contactNo must not exceed 12 digits")
	}
	return nil
}

// ValidateEmail checks that the email address matches a standard format.
func ValidateEmail(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("email must not be empty")
	}
	if !emailRegex.MatchString(value) {
		return fmt.Errorf("email is not a valid email address")
	}
	return nil
}

// ValidateEmployeeCreate performs all field-level validations for creating an employee.
// Returns the first validation error encountered, or nil if all fields are valid.
func ValidateEmployeeCreate(e *model.EmployeeCreate) error {
	if err := ValidateEmpID(e.EmpID); err != nil {
		return err
	}
	if err := ValidateString("name", e.Name); err != nil {
		return err
	}
	if err := ValidateEmail(e.Email); err != nil {
		return err
	}
	if err := ValidateContactNo(e.ContactNo); err != nil {
		return err
	}
	if err := ValidateString("role", e.Role); err != nil {
		return err
	}
	if err := ValidateString("department", e.Department); err != nil {
		return err
	}
	if e.Salary < 0 {
		return fmt.Errorf("salary must be greater than or equal to 0")
	}
	if err := ValidateString("created_by", e.CreatedBy); err != nil {
		return err
	}
	return nil
}

// ValidateEmployeeUpdate performs field-level validations for updating an employee.
// Only non-nil fields are validated. Returns an error if any provided field is invalid.
func ValidateEmployeeUpdate(e *model.EmployeeUpdate) error {
	if e.Name != nil {
		if err := ValidateString("name", *e.Name); err != nil {
			return err
		}
	}
	if e.Email != nil {
		if err := ValidateEmail(*e.Email); err != nil {
			return err
		}
	}
	if e.ContactNo != nil {
		if err := ValidateContactNo(*e.ContactNo); err != nil {
			return err
		}
	}
	if e.Role != nil {
		if err := ValidateString("role", *e.Role); err != nil {
			return err
		}
	}
	if e.Department != nil {
		if err := ValidateString("department", *e.Department); err != nil {
			return err
		}
	}
	if e.Salary != nil && *e.Salary < 0 {
		return fmt.Errorf("salary must be greater than or equal to 0")
	}
	return nil
}
