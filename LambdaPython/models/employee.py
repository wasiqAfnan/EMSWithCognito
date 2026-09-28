from pydantic import BaseModel, EmailStr, Field, field_validator

from utils.validators import (
    validate_string,
    validate_emp_id,
    validate_contact_no,
)


class EmployeeCreate(BaseModel):
    empId: str = Field(min_length=6)
    name: str = Field(min_length=1)
    email: EmailStr
    contactNo: str
    role: str = Field(min_length=1)
    department: str = Field(min_length=1)
    salary: float = Field(ge=0)
    created_by: str = Field(min_length=1)

    # String validation
    _validate_strings = field_validator(
        "name", "role", "department", "created_by",
        mode="before"
    )(validate_string)


    # Employee ID validation
    _validate_emp_id = field_validator(
        "empId",
        mode="before"
    )(validate_emp_id)

    # Contact number validation
    _validate_contact_no = field_validator(
        "contactNo",
        mode="before"
    )(validate_contact_no)


class EmployeeUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1)
    email: EmailStr | None = None
    contactNo: str | None = Field(
        default=None,
        min_length=10,
        max_length=12
    )
    role: str | None = Field(default=None, min_length=1)
    department: str | None = Field(default=None, min_length=1)
    salary: float | None = Field(default=None, ge=0)

    # String validation
    _validate_name = field_validator(
        "name",
        mode="before"
    )(validate_string)

    _validate_role = field_validator(
        "role",
        mode="before"
    )(validate_string)

    _validate_department = field_validator(
        "department",
        mode="before"
    )(validate_string)

    # Contact number validation
    _validate_contact_no = field_validator(
        "contactNo",
        mode="before"
    )(validate_contact_no)