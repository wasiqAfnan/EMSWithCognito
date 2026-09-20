# IMPORTANT! Line no 2 to 15 is for connecting to MongoDB using .env file
import os
from pymongo import MongoClient
from dotenv import load_dotenv

load_dotenv()

MONGODB_URI = os.environ["MONGODB_URI"]
MONGODB_DATABASE = os.environ["MONGODB_DATABASE"]

client = MongoClient(MONGODB_URI)

db = client[MONGODB_DATABASE]

employees_collection = db["employees"]
users_collection = db["users"]


# IMPORTANT! Line No 19 to 31 is for connecting to MongoDB using boto3 package and parameter store in aws lambda
# from pymongo import MongoClient

# from utils.parameter_store import get_parameter


# MONGODB_URI = get_parameter("/ems/prod/mongodb-uri")
# MONGODB_DATABASE = get_parameter("/ems/prod/mongodb-database")

# client = MongoClient(MONGODB_URI)

# db = client[MONGODB_DATABASE]

# employees_collection = db["employees"]
# users_collection = db["users"]