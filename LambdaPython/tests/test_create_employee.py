import json

from lambdas.employees.create_employee import lambda_handler


event = {
    "body": json.dumps({
        "empId": "EMP001",
        "name": "Wasiq Afnan Ansari",
        "email": "wasiq@example.com",
        "role": "Cloud Developer",
        "department": "Cloud",
        "salary": 700000
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))