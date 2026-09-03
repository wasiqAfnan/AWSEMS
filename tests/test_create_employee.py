import json

from lambdas.employees.create_employee import lambda_handler


event = {
    "body": json.dumps({
        "empId": "EMP003",
        "name": "Abhijeet",
        "email": "abhijeet@gmail.com",
        "contactNo": "1234567890",
        "role": "Cloud Engineer",
        "department": "Cloud",
        "salary": 400000
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))