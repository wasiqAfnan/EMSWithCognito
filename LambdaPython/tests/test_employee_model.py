from models.employee import EmployeeCreate, EmployeeUpdate


employee = EmployeeCreate(
    empId="EMP001",
    name="Rahul Kumar",
    email="rahul@example.com",
    role="Software Engineer",
    department="Engineering",
    salary=60000
)

print(employee)
print(employee.model_dump())


update = EmployeeUpdate(
    salary=65000
)

print(update)
print(update.model_dump(exclude_none=True))