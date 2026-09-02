from pydantic import BaseModel, EmailStr, Field


class EmployeeCreate(BaseModel):
    empId: str = Field(min_length=1)
    name: str = Field(min_length=1)
    email: EmailStr
    contactNo: str = Field(min_length=1, max_length=10)
    role: str = Field(min_length=1)
    department: str = Field(min_length=1)
    salary: float = Field(ge=0)


class EmployeeUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1)
    email: EmailStr | None = None
    contactNo: str | None = Field(default=None, min_length=1)
    role: str | None = Field(default=None, min_length=1)
    department: str | None = Field(default=None, min_length=1)
    salary: float | None = Field(default=None, ge=0)