import json

from lambdas.employees.update_employee import lambda_handler


event = {
    "pathParameters": {
        "empId": "EMP001"
    },
    "body": json.dumps({
        "name": "Wasiq Afnan Ansari",
        "email": "wasiq@example.com",
        "role": "Cloud Developer",
        "department": "Cloud",
        "salary": 700000
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))