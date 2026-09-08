import json

from lambdas.employees.create_employee import lambda_handler


event = {
    "body": json.dumps({
        "empId": "EMP005",
        "name": "Nilarpan",
        "email": "nilarpan2@gmail.com",
        "contactNo": "12345678091",
        "role": "Database Admin",
        "department": "Database",
        "salary": 500000
    })
}

response = lambda_handler(event, None)

print(json.dumps(response, indent=4))