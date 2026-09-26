import json

from lambdas.employees.get_employees import lambda_handler


event = {}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))