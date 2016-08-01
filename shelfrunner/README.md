To have Iron.io run this example job:

1. Modify the Makefile to upload to the Docker repository of your choice.
2. Run `make`.
3. Follow the Iron.io documentation to start an API server and at least one Runner.
4. Start the job via (assuming Iron.io and Mongo is running on localhost):

    ```sh
    curl -X POST -H "Content-Type: application/json" -d '{
        "jobs": [{
        "image": "dwhitena/shelfrunner",
        "payload": "{\"asset\": \"579f803559b0b30009ff60f1\", \"mongo_host\": \"localhost\"}"
      }]
    }' "http://localhost:8080/v1/groups/mygroup/jobs"
```
