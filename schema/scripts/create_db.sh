if [ -z "$1" ]; then
    echo "Usage: create_db.sh <database.json>"
    exit 1
fi

source schema/scripts/token.env

curl -u $WDB_USERNAME:$WDB_PASSWORD \
    --location "$BASE_URL/api/databases" \
    --header "Content-Type: application/json" \
    --data @$1 | jq .
