FROM mongo

ADD example-items.json /
CMD mongoimport --host mongodb --db coral-poc --collection Items --type json --file /example-items.json --jsonArray
