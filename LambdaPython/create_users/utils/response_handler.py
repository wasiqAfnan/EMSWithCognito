from pydantic import ValidationError

from models.response import APIResponse
import json


def success_response(status_code, message, data=None):
    response = APIResponse(
        statusCode=status_code,
        message=message,
        data=data
    )

    return {
        "success": True,
        "statusCode": status_code,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps(response.model_dump())
    }


def error_response(status_code, message, data=None):
    response = APIResponse(
        statusCode=status_code,
        message=message,
        data=data
    )

    return {
        "success": False,
        "statusCode": status_code,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps(response.model_dump())
    }


def validation_error_response(error: ValidationError):
    first_error = error.errors()[0]

    field = first_error["loc"][0]
    message = first_error["msg"]

    if message.startswith("Value error, "):
        message = message.replace("Value error, ", "", 1)

    return error_response(
        400,
        f"Invalid employee data. {field} {message}"
    )
