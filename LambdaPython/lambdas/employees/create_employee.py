import json

from pydantic import ValidationError
from pymongo.errors import DuplicateKeyError

from db.mongodb import employees_collection
from models.employee import EmployeeCreate
from utils.response_handler import (
    success_response,
    error_response,
    validation_error_response,
)


def lambda_handler(event, context):
    try:
        # Extract request body
        raw_body = event.get("body")
        print(f"[Lambda create_employee] Received raw_body from API Gateway: {raw_body}")

        if not raw_body:
            return error_response(
                400,
                "Request body is required"
            )

        # Parse JSON body
        body = json.loads(raw_body)
        print(f"[Lambda create_employee] Parsed JSON body: {body}")

        # Validate request data
        employee = EmployeeCreate(**body)

        # Convert validated Pydantic model to dictionary
        employee_data = employee.model_dump()
        print(f"[Lambda create_employee] Validated employee_data being inserted to MongoDB: {employee_data}")

        # Check if Employee ID already exists
        if employees_collection.find_one({"empId": employee.empId}):
            return error_response(
                409,
                "Employee ID already exists"
            )

        # Check if email already exists
        if employees_collection.find_one({"email": str(employee.email)}):
            return error_response(
                409,
                "Email already exists"
            )

        # Insert employee into MongoDB
        result = employees_collection.insert_one(employee_data)

        # Convert MongoDB ObjectId to string
        employee_data["_id"] = str(result.inserted_id)

        return success_response(
            201,
            "Employee created successfully",
            employee_data
        )

    except json.JSONDecodeError:
        return error_response(
            400,
            "Invalid JSON body"
        )

    except ValidationError as e:
        return validation_error_response(e)

    except DuplicateKeyError as e:
        error_message = str(e)

        if "empId" in error_message:
            message = "Employee ID already exists"

        elif "email" in error_message:
            message = "Email already exists"


        else:
            message = "Employee with the same unique field already exists"

        return error_response(
            409,
            message
        )

    except Exception as e:
        print(f"Error creating employee: {e}")

        return error_response(
            500,
            "Internal server error"
        )