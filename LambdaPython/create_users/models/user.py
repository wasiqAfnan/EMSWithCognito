from pydantic import BaseModel, EmailStr, Field, field_validator

from utils.validators import validate_string


class UserCreate(BaseModel):
    cognito_sub: str = Field(min_length=1)
    email: EmailStr
    name: str = Field(min_length=1)
    role: str = Field(min_length=1)

    # String validation
    _validate_strings = field_validator(
        "cognito_sub", "name", "role",
        mode="before"
    )(validate_string)
