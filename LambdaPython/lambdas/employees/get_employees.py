from db.mongodb import employees_collection
from utils.response_handler import (
    success_response,
    error_response,
)
import json


def lambda_handler(event, context):
    try:
        # Extract created_by from JSON body
        raw_body = event.get("body")
        created_by = None
        if raw_body:
            body_data = json.loads(raw_body)
            created_by = body_data.get("created_by")

        if not created_by:
            return error_response(400, "Missing created_by parameter")

        # Fetch all employees from MongoDB
        employees = list(
            employees_collection.find(
                {"created_by": created_by}, # return filtered employees
                {
                    "_id": 0 # exclude _id field
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