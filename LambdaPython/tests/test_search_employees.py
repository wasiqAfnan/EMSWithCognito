import json

from lambdas.employees.search_employee import lambda_handler

event = {"queryStringParameters": {"q": "wasiq@example.com"}}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))
