import json
from os import environ
import os
import requests


BASE_URL = environ.get("BASE_URL")
WDB_USERNAME = environ.get("WDB_USERNAME")
WDB_PASSWORD = environ.get("WDB_PASSWORD")

print(WDB_PASSWORD, WDB_USERNAME, BASE_URL)


# get the directory where this file is located
dir_path = os.path.dirname(os.path.realpath(__file__))
data_dir_path = f"{dir_path}/../data"

# iterate over directories in data directory
for dir in os.listdir(data_dir_path):
    DATABASE_NAME: str
    COLLECTION_NAME: str
    EXEMPTED_SCHEMA: bool

    with open(f"{data_dir_path}/{dir}/collection.json") as f:
        data = json.loads(f.read())
        DATABASE_NAME = data["database"]
        COLLECTION_NAME = data["collection"]
        EXEMPTED_SCHEMA = data["exempt"]

    if EXEMPTED_SCHEMA:
        print(f"Skipping {COLLECTION_NAME} as it is exempted")
        continue

    url = f"{BASE_URL}/api/databases/{DATABASE_NAME}/collections/{COLLECTION_NAME}"

    payload = {}
    headers = {}
    response = requests.request(
        "GET", url, headers=headers, data=payload, auth=(WDB_USERNAME, WDB_PASSWORD)
    )

    response_json = json.loads(response.text)
    collection_records = response_json["response"]["records"]

    kPresent = True if "Kafka" in collection_records else False
    print(f"Collection {COLLECTION_NAME} has  Kafka : {kPresent}")
