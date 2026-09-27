from db.mongodb import employees_collection
from utils.response_handler import (
    success_response,
    error_response,
)
from utils.validators import validate_emp_id
import json


def lambda_handler(event, context):
    try:
        # Extract Employee ID from path parameters
        path_parameters = event.get("pathParameters") or {}
        emp_id = path_parameters.get("empId")

        if not emp_id:
            return error_response(
                400,
                "Employee ID is required"
            )

        # Validate Employee ID
        emp_id = validate_emp_id(emp_id)

        # Extract created_by from JSON body
        raw_body = event.get("body")
        created_by = None
        if raw_body:
            body_data = json.loads(raw_body)
            created_by = body_data.get("created_by")

        if not created_by:
            return error_response(400, "Missing created_by parameter")

        # Check if employee exists and belongs to user
        employee = employees_collection.find_one(
            {"empId": emp_id, "created_by": created_by}
        )

        if not employee:
            return error_response(
                404,
                "Employee not found"
            )

        # Delete employee
        employees_collection.delete_one(
            {"empId": emp_id, "created_by": created_by}
        )

        return success_response(
            200,
            "Employee deleted successfully",
            {
                "empId": emp_id
            }
        )

    except ValueError as e:
        return error_response(
            400,
            f"Invalid employee ID. {str(e)}"
        )

    except Exception as e:
        print(f"Error deleting employee: {e}")

        return error_response(
            500,
            "Internal server error"
        )