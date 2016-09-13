    sponged
    xeniad
    sponge item upsert -p input/items.json
    xenia pattern upsert -p metadata/patterns/
    xenia relationship upsert -p metadata/relationships/
    xenia view upsert -p metadata/views/
    wire graph add -p input/items/
    wire view execute -n VTEST_thread -i c1b2bbfe-af9f-4903-8777-bd47c4d5b20a
