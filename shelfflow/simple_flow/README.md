# Simple flow of items from input to view.

1. Declare environmental variables required by the `sponged`, `sponge`, `xeniad`, `xenia`, and `wire` commands.

2. Run `sponged` and `xeniad`.

3. Run `runflow.sh`.  This will import the items, import the metadata, and build the graph.

4. Now you can examine enabled views via, e.g.,

    ```
    wire view execute -n thread -i c1b2bbfe-af9f-4903-8777-bd47c4d5b20a
    wire view execute -n "user comments" -i a63af637-58af-472b-98c7-f5c00743bac6
    wire view execute -n "user comments" -i 80aa936a-f618-4234-a7be-df59a14cf8de
    ```

5. Cleanup by running `cleanup.sh`.
    
