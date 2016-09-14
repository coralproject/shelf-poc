#!/bin/bash

# Remove the graph quads.
wire graph remove -p cleanup/quadparams/quadparams.json

# Remove the items from mongo.
for line in $(cat "cleanup/items/items.txt"); do
	sponge item delete -i "$line"
done

# Remove the patterns from mongo.
for line in $(cat "cleanup/patterns/patterns.txt"); do
	xenia pattern delete -t "$line"
done

# Remove the relationships from mongo.
for line in $(cat "cleanup/relationships/relationships.txt"); do
	xenia relationship delete -p "$line"
done

# Remove the views from mongo.
for line in $(cat "cleanup/views/views.txt"); do
	xenia view delete -n "$line"
done
