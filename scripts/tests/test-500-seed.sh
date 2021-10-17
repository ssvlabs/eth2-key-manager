#!/bin/bash

current_folder=$(pwd)
scripts_folder=$current_folder/scripts/tests

cd $current_folder/cli/
if [ ! -f  "./cli" ] ; then
    echo 'Compiling fresh eth cli..'
    go build
fi
cd $current_folder

accounts_number=499
highest_start_from=123456
highest=$($scripts_folder/get-highest-values.py $accounts_number $highest_start_from)

echo 'Generating number of accounts:' "$accounts_number"

echo 'Generating seed accounts'
$current_folder/cli/cli wallet account create \
  --seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff \
  --index=$accounts_number \
  --accumulate \
	--response-type=object \
	--highest-source=$highest \
	--highest-target=$highest \
	--highest-proposal=$highest \
	--network=prater \
	> $scripts_folder/seed-accounts-$accounts_number.json

echo 'Generating seedless accounts'
private_keys=$($scripts_folder/get-private-keys.py $accounts_number)

index_from=20
highest_accounts=$(($accounts_number + $index_from))
highest=$($scripts_folder/get-highest-values.py $highest_accounts $highest_start_from)
$current_folder/cli/cli wallet account create-seedless \
  --highest-proposal=$highest \
  --highest-source=$highest \
  --highest-target=$highest \
  --network=pyrmont \
  --index-from=$index_from \
  --response-type=object \
  --private-key=$private_keys \
  > $scripts_folder/seedless-accounts-$accounts_number.json
