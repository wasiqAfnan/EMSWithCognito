import json

from lambdas.users.get_user import lambda_handler


event = {
    "pathParameters": {
        "id": "test-uuid-1234-5678"
    }
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))

