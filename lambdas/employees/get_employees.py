from db.mongodb import employees_collection
from utils.response_handler import (
    success_response,
    error_response,
)


def lambda_handler(event, context):
    try:
        # Fetch all employees from MongoDB
        employees = list(
            employees_collection.find(
                {},
                {
                    "_id": 0
                }
            )
        )

        return success_response(
            200,
            "Employees fetched successfully",
            employees
        )

    except Exception as e:
        print(f"Error fetching employees: {e}")

        return error_response(
            500,
            "Internal server error"
        )