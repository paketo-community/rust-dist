#!/bin/sh
# update the versions of rust being used in buildpack.toml
#   requires `jq` and `yj` to be on the PATH
set -e

# pull all the available versions of rust from the dep server
#   - pick the first & second ones (assumes they are sorted most recent first)
#   - remaps the name field to the id field and adds a name of "Rust"
#   - removes some currently unused fields (.cpe,.created_at,.modified_at,.deprecation_date)
#   - reformats stacks
DEPS=$(curl -s "https://api.deps.paketo.io/v1/dependency?name=rust" | jq '[first, nth(1)] | map({"id": .name} + . + {"name": "Rust"}) | map(del(.cpe,.created_at,.modified_at,.deprecation_date)) | map(. + {"stacks": [.stacks[].id]})')

# Update buildpack.toml with the two most recent deps
#   - delete the old deps section
#   - add in the new one, uses --jsonargs to feed input from previous command
#   - reformat to toml & update buildpack.toml
OLD=$(yj -t < buildpack.toml)
echo "$OLD" | jq --argjson DEPS "$DEPS" 'del(.metadata.dependencies) + {"metadata": (.metadata + {"dependencies": $DEPS})}' | yj -i -jt > buildpack.toml
