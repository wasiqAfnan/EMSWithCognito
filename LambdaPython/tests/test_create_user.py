import json

from lambdas.users.create_user import lambda_handler


event = {
    "body": json.dumps({
        "cognito_sub": "test-uuid-1234-5678",
        "name": "Jane Doe",
        "email": "jane.doe@example.com",
        "role": "USER"
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))
