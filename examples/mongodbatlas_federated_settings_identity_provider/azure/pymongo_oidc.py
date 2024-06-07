"""This script is meant as a minimal example of how to connect to MongoDB using https://pymongo.readthedocs.io/en/stable/index.html.
Note how there is nothing specific about using OIDC, it only requires a valid `MONGODB_URI`, `DATABASE`, and `COLLECTION`.
See https://www.mongodb.com/docs/manual/reference/connection-string/#connection-string-formats for how to specify `MONGODB_URI`.

Feel free to replace this with your awesome python script!"""
from datetime import datetime
from json import loads
from os import getenv

from pymongo import MongoClient


def insert_record():
    uri = getenv("MONGODB_URI")
    assert uri, "missing MONGODB_URI"
    print(f"creating client with uri={uri}")
    client = MongoClient(uri)

    database = getenv("DATABASE")
    collection = getenv("COLLECTION")
    err_msg = f"missing database/collection {database}/{collection}"
    assert database and collection, err_msg
    if record_str := getenv("RECORD"):
        record = loads(record_str)
    else:
        record = {"hello": "world"}
    record["ts"] = datetime.now().isoformat()
    print(f"inserting into {database} {collection}, record: {record}")
    response = (
        client.get_database(database).get_collection(collection).insert_one(record)
    )
    print(f"insert response: {response}")
    client.close()
    print("script complete")


if __name__ == "__main__":
    insert_record()
