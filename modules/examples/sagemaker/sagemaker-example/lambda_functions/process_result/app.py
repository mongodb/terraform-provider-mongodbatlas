import json
import pymongo
import os


def handler(event, context):
    """
    Lambda handler. Writes results back to MongoDB.
    """

    print(json.dumps(event))

    try:
        # Read environment variables.
        MONGO_ENDPOINT = os.environ['mongo_endpoint']
        MONGO_DB = os.environ['dbname']
        MONGO_COL = "predictions"

        # Connect to MongoDB Atlas.
        client = pymongo.MongoClient(MONGO_ENDPOINT)
        db = client[MONGO_DB]
        col = db[MONGO_COL]

        # Insert result.
        result = col.insert_one(event['detail'])
        print("Inserted: ".format(result.inserted_id))

        return "Inserted: ".format(result.inserted_id)
    except Exception as ex:
        print("Exception: " + str(ex))
        raise ex
