import json

from pydantic import ValidationError
from pymongo.errors import DuplicateKeyError

from db.mongodb import users_collection
from models.user import UserCreate
from utils.response_handler import (
    success_response,
    error_response,
    validation_error_response,
)


def lambda_handler(event, context):
    try:
        # Extract request body
        raw_body = event.get("body")

        if not raw_body:
            return error_response(
                400,
                "Request body is required"
            )

        # Parse JSON body
        body = json.loads(raw_body)

        # Validate request data
        user = UserCreate(**body)

        # Convert validated Pydantic model to dictionary
        user_data = user.model_dump()

        # Check if email already exists
        if users_collection.find_one({"email": str(user.email)}):
            return error_response(
                409,
                "Email already exists"
            )
            
        # Check if cognito_sub already exists
        if users_collection.find_one({"cognito_sub": user.cognito_sub}):
            return error_response(
                409,
                "User with this Cognito Sub already exists"
            )

        # Insert user into MongoDB
        result = users_collection.insert_one(user_data)

        # Convert MongoDB ObjectId to string
        user_data["_id"] = str(result.inserted_id)

        return success_response(
            201,
            "User created successfully",
            user_data
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
        elif "cognito_sub" in error_message:
            message = "Cognito Sub already exists"
        else:
            message = "User with the same unique field already exists"

        return error_response(
            409,
            message
        )

    except Exception as e:
        print(f"Error creating user: {e}")

        return error_response(
            500,
            "Internal server error"
        )
