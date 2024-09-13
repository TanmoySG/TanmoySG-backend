import json
import os
import requests

BASE_URL = os.environ.get("BASE_URL")
WDB_USERNAME = os.environ.get("WDB_USERNAME")
WDB_PASSWORD = os.environ.get("WDB_PASSWORD")


def patch(
    base_url, database_name, collection_name, payload, pKeyField, pKeyValue
) -> str:
    url = f"{base_url}/api/databases/{database_name}/collections/{collection_name}/records?key={pKeyField}&value={pKeyValue}"

    p = json.dumps(payload)
    headers = {
        "Content-Type": "application/json",
    }

    response = requests.request(
        "PATCH", url, headers=headers, data=p, auth=(WDB_USERNAME, WDB_PASSWORD)
    )

    return response.text


def create(base_url, database_name, collection_name, payload) -> str:
    url = f"{base_url}/api/databases/{database_name}/collections/{collection_name}/records"

    p = json.dumps(payload)
    headers = {
        "Content-Type": "application/json",
    }

    response = requests.request(
        "POST", url, headers=headers, data=p, auth=(WDB_USERNAME, WDB_PASSWORD)
    )

    return response.text


# get the directory where app.py file is located
dir_path = os.path.dirname(os.path.realpath(__file__))
data_dir_path = f"{dir_path}/../data"

# iterate over directories in data directory
for dir in os.listdir(data_dir_path):
    DATABASE_NAME: str
    COLLECTION_NAME: str
    EXEMPTED_SCHEMA: bool
    RECORDS_ARRAY: list

    with open(f"{data_dir_path}/{dir}/collection.json") as f:
        collection_config = json.loads(f.read())
        DATABASE_NAME = collection_config["database"]
        COLLECTION_NAME = collection_config["collection"]
        EXEMPTED_SCHEMA = collection_config["exempt"]

    if EXEMPTED_SCHEMA:
        print(f"Skipping {COLLECTION_NAME} as it is exempted")
        continue

    with open(f"{data_dir_path}/{dir}/records.json") as f:
        RECORDS_ARRAY = json.loads(f.read())

    url = f"{BASE_URL}/api/databases/{DATABASE_NAME}/collections/{COLLECTION_NAME}"

    payload = {}
    headers = {}
    response = requests.request(
        "GET", url, headers=headers, data=payload, auth=(WDB_USERNAME, WDB_PASSWORD)
    )

    response_json = json.loads(response.text)
    collection_records = response_json["response"]["records"]
    primary_key_field = response_json["response"]["primaryKey"]

    for record in RECORDS_ARRAY:
        if record[primary_key_field] not in collection_records:
            res = create(
                BASE_URL,
                DATABASE_NAME,
                COLLECTION_NAME,
                record,
            )
            print(
                f"Creating new record for {primary_key_field}={record[primary_key_field]} in {COLLECTION_NAME}. Response: {res}"
            )
        else:
            res = patch(
                BASE_URL,
                DATABASE_NAME,
                COLLECTION_NAME,
                record,
                primary_key_field,
                record[primary_key_field],
            )
            print(
                f"Updating existing record for {primary_key_field}={record[primary_key_field]} in {COLLECTION_NAME}. Response: {res}"
            )
