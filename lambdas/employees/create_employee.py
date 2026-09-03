import json

from pymongo.errors import DuplicateKeyError
from pydantic import ValidationError

from db.mongodb import employees_collection
from models.employee import EmployeeCreate
from models.response import APIResponse


def lambda_handler(event, context):
    try:
        # Extract request body
        raw_body = event.get("body")

        if not raw_body:
            return APIResponse(
                statusCode=400,
                message="Request body is required",
                data=None
            ).model_dump()

        # Parse JSON body
        body = json.loads(raw_body)

        # Validate request data
        employee = EmployeeCreate(**body)

        # Convert validated Pydantic model to dictionary
        employee_data = employee.model_dump()

        # Check if Employee ID already exists
        if employees_collection.find_one({"empId": employee.empId}):
            return APIResponse(
                statusCode=409,
                message="Employee ID already exists",
                data=None
            ).model_dump()

        # Check if email already exists
        if employees_collection.find_one({"email": str(employee.email)}):
            return APIResponse(
                statusCode=409,
                message="Email already exists",
                data=None
            ).model_dump()

        # Check if contact number already exists
        if employees_collection.find_one({"contactNo": employee.contactNo}):
            return APIResponse(
                statusCode=409,
                message="Contact number already exists",
                data=None
            ).model_dump()

        # Insert employee into MongoDB
        result = employees_collection.insert_one(employee_data)

        # Add MongoDB generated ID to response
        employee_data["_id"] = str(result.inserted_id)

        return APIResponse(
            statusCode=201,
            message="Employee created successfully",
            data=employee_data
        ).model_dump()

    except json.JSONDecodeError:
        return APIResponse(
            statusCode=400,
            message="Invalid JSON body",
            data=None
        ).model_dump()

    except ValidationError as e:
        errors = []

        for error in e.errors():
            errors.append({
                "field": ".".join(str(location) for location in error["loc"]),
                "message": error["msg"],
            })

        return APIResponse(
            statusCode=400,
            message="Invalid employee data",
            data=errors
        ).model_dump()

    except DuplicateKeyError as e:
        # Determine which unique field caused the conflict
        error_message = str(e)

        if "empId" in error_message:
            message = "Employee ID already exists"

        elif "email" in error_message:
            message = "Email already exists"

        elif "contactNo" in error_message:
            message = "Contact number already exists"

        else:
            message = "Employee with the same unique field already exists"

        return APIResponse(
            statusCode=409,
            message=message,
            data=None
        ).model_dump()

    except Exception as e:
        print(f"Error creating employee: {e}")

        return APIResponse(
            statusCode=500,
            message="Internal server error",
            data=None
        ).model_dump()