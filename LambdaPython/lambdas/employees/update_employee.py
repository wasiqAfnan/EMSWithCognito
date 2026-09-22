import json

from pydantic import ValidationError
from pymongo.errors import DuplicateKeyError

from db.mongodb import employees_collection
from models.employee import EmployeeUpdate
from utils.response_handler import (
    success_response,
    error_response,
    validation_error_response,
)


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

        # Extract request body
        raw_body = event.get("body")

        if not raw_body:
            return error_response(
                400,
                "Request body is required"
            )

        # Parse JSON body
        body = json.loads(raw_body)

        # Validate update data using Pydantic
        employee_update = EmployeeUpdate(**body)

        # Convert validated data to dictionary
        update_data = employee_update.model_dump(
            exclude_none=True
        )

        # Make sure at least one field is provided
        if not update_data:
            return error_response(
                400,
                "At least one field is required for update"
            )

        # Extract created_by from the injected body
        created_by = body.get("created_by")

        if not created_by:
            return error_response(400, "Missing created_by parameter")

        # Check if employee exists and belongs to the user
        existing_employee = employees_collection.find_one(
            {"empId": emp_id, "created_by": created_by}
        )

        if not existing_employee:
            return error_response(
                404,
                "Employee not found"
            )

        # Check if email already belongs to another employee
        if "email" in update_data:
            existing_email = employees_collection.find_one({
                "email": str(update_data["email"]),
                "empId": {"$ne": emp_id}
            })

            if existing_email:
                return error_response(
                    409,
                    "Email already exists"
                )

            update_data["email"] = str(update_data["email"])

        # Check if contact number already belongs to another employee
        if "contactNo" in update_data:
            existing_contact = employees_collection.find_one({
                "contactNo": update_data["contactNo"],
                "empId": {"$ne": emp_id}
            })

            if existing_contact:
                return error_response(
                    409,
                    "Contact number already exists"
                )

        # Update employee
        employees_collection.update_one(
            {"empId": emp_id, "created_by": created_by},
            {"$set": update_data} # update only the specified fields
        )

        # Fetch updated employee
        updated_employee = employees_collection.find_one(
            {"empId": emp_id, "created_by": created_by},
            {"_id": 0}
        )

        return success_response(
            200,
            "Employee updated successfully",
            updated_employee
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

        if "email" in error_message:
            message = "Email already exists"

        elif "contactNo" in error_message:
            message = "Contact number already exists"

        else:
            message = "Employee with the same unique field already exists"

        return error_response(
            409,
            message
        )

    except Exception as e:
        print(f"Error updating employee: {e}")

        return error_response(
            500,
            "Internal server error"
        )