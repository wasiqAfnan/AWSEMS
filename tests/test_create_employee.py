import json

from lambdas.employees.create_employee import lambda_handler


event = {
    "body": json.dumps({
        "empId": "EMP004",
        "name": "Nilarpan",
        "email": "nilarpan@gmail.com",
        "contactNo": "1234567891",
        "role": "Database Admin",
        "department": "Database",
        "salary": 500000
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))