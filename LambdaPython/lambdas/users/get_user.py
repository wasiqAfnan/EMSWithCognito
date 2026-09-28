from db.mongodb import users_collection
from utils.response_handler import (
    success_response,
    error_response,
)


def lambda_handler(event, context):
    try:
        # Extract cognito_sub from path parameters
        path_parameters = event.get("pathParameters") or {}
        cognito_sub = path_parameters.get("id")

        if not cognito_sub:
            return error_response(
                400,
                "Cognito Sub is required"
            )

        # Fetch user from MongoDB
        user = users_collection.find_one(
            {"cognito_sub": cognito_sub}
        )

        if not user:
            return error_response(
                404,
                "User not found"
            )

        # Convert MongoDB ObjectId to string
        user["_id"] = str(user["_id"])

        return success_response(
            200,
            "User fetched successfully",
            user
        )

    except Exception as e:
        print(f"Error fetching user: {e}")

        return error_response(
            500,
            "Internal server error"
        )
