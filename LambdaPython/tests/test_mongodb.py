from db.mongodb import client, employees_collection

try:
    client.admin.command("ping")
    print("MongoDB connection successful!")

    print(f"Database: {employees_collection.database.name}")
    print(f"Collection: {employees_collection.name}")

except Exception as e:
    print("MongoDB connection failed!")
    print(e)