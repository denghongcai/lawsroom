#!/bin/bash

for name in $(ls src | grep -P "\.md" | cut -d. -f 1)
do
    pandoc --from=markdown --to=rst --output=source/${name}.rst src/$name.md
done

