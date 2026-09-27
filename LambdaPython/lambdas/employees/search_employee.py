from db.mongodb import employees_collection
from utils.response_handler import (
    success_response,
    error_response,
)

import json


def lambda_handler(event, context):
    try:
        # Get search query from query string parameters
        query_parameters = event.get("queryStringParameters") or {}
        search = query_parameters.get("q", "").strip()

        if not search:
            return error_response(400, "Search query is required")

        # Extract created_by from JSON body
        raw_body = event.get("body")
        created_by = None
        if raw_body:
            body_data = json.loads(raw_body)
            created_by = body_data.get("created_by")
        if not created_by:
            return error_response(400, "Missing created_by parameter")

        # Search across multiple employee fields
        search_filter = {
            "created_by": created_by,
            "$or": [
                {"empId": {"$regex": search, "$options": "i"}},
                {"email": {"$regex": search, "$options": "i"}},
                {"name": {"$regex": search, "$options": "i"}},
                {"role": {"$regex": search, "$options": "i"}},
                {"department": {"$regex": search, "$options": "i"}},
            ]
        }

        # Fetch matching employees
        employees = list(employees_collection.find(search_filter, {"_id": 0}))

        return success_response(
            200, "Employee search completed successfully", employees
        )

    except Exception as e:
        print(f"Error searching employees: {e}")

        return error_response(500, "Internal server error")
