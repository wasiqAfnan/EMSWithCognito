import json

from lambdas.employees.delete_employee import lambda_handler


event = {
    "pathParameters": {
        "empId": "EMP005"
    }
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))