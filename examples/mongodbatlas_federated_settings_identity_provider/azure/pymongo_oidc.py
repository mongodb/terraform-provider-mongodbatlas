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
