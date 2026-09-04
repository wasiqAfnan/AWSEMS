from pydantic import BaseModel, EmailStr, Field, field_validator


class EmployeeCreate(BaseModel):
    empId: str = Field(min_length=6)
    name: str = Field(min_length=1)
    email: EmailStr
    contactNo: str = Field()
    role: str = Field(min_length=1)
    department: str = Field(min_length=1)
    salary: float = Field(ge=0)

    @field_validator(
        "empId",
        "name",
        "contactNo",
        "role",
        "department",
        mode="before"
    )
    @classmethod
    def validate_strings(cls, value):
        if not isinstance(value, str):
            raise ValueError("must be a string")

        value = value.strip()

        if not value:
            raise ValueError("must not be empty")

        return value

    @field_validator("empId")
    @classmethod
    def validate_emp_id(cls, value):
        if not value.startswith("EMP"):
            raise ValueError("must start with EMP")

        if not value[3:].isdigit():
            raise ValueError("must contain only digits after EMP")

        if len(value) < 6:
            raise ValueError("must be at least 6 characters long")

        return value

    @field_validator("contactNo")
    @classmethod
    def validate_contact_no(cls, value):
        if not value.isdigit():
            raise ValueError("must contain only digits")

        if len(value) < 10:
            raise ValueError("must be atleast 10 digits")

        if len(value) > 12:
                    raise ValueError("must not exceeds 12 digits")
        return value


class EmployeeUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1)
    email: EmailStr | None = None
    contactNo: str | None = Field(default=None, min_length=10, max_length=12)
    role: str | None = Field(default=None, min_length=1)
    department: str | None = Field(default=None, min_length=1)
    salary: float | None = Field(default=None, ge=0)

    @field_validator(
        "name",
        "contactNo",
        "role",
        "department",
        mode="before"
    )
    @classmethod
    def strip_and_validate_optional_strings(cls, value):
        if value is None:
            return None

        if not isinstance(value, str):
            raise ValueError("Must be a string")

        value = value.strip()

        if not value:
            raise ValueError("Must not be empty or contain only whitespace")

        return value