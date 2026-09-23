package model

// EmployeeCreate holds the fields required to create a new employee.
type EmployeeCreate struct {
	EmpID      string  `json:"empId"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	ContactNo  string  `json:"contactNo"`
	Role       string  `json:"role"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
	CreatedBy  string  `json:"created_by"` // Cognito sub of the owner
}

// EmployeeUpdate holds the optional fields that can be updated for an employee.
// Pointer types allow distinguishing between a missing field and a zero value.
type EmployeeUpdate struct {
	Name       *string  `json:"name,omitempty"`
	Email      *string  `json:"email,omitempty"`
	ContactNo  *string  `json:"contactNo,omitempty"`
	Role       *string  `json:"role,omitempty"`
	Department *string  `json:"department,omitempty"`
	Salary     *float64 `json:"salary,omitempty"`
}
