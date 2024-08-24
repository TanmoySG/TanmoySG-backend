# for each directory find file with extention .collection.json

collections_dir="schema/databases/tsg-on-web_v0_beta/collections"

for dir in $collections_dir/*/     # list directories in the form "/tmp/dirname/"
do
    dir=${dir%*/}      # remove the trailing "/"
    sub_dir="${dir##*/}"    # print everything after the final "/"

    sh schema/scripts/create_collection.sh tsg-on-web_v0_beta_1 $collections_dir/$sub_dir/$sub_dir.collection.json
done