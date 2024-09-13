# check if $1 and `$2` are passed as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: create_collection.sh <database_name> <collection.json>"
    exit 1
fi

source schema/scripts/token.env

schema_fp=$(echo $2 | sed 's/\(.*\)collection/\1schema/')

jq '.schema = input' $2 $schema_fp > $2.tmp

curl -X POST -u $WDB_USERNAME:$WDB_PASSWORD \
    --location "$BASE_URL/api/databases/$1/collections" \
    --header 'Content-Type: application/json' \
    --data @$2.tmp | jq .

rm $2.tmp