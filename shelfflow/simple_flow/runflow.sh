#!/bin/bash

# Import the metadata required to generate a view.
xenia pattern upsert -p metadata/patterns/
xenia relationship upsert -p metadata/relationships/
xenia view upsert -p metadata/views/

# Import the items.
sponge item upsert -p input/items/


