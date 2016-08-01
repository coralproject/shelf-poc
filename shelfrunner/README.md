To have Iron.io run this example job:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
    "jobs": [{
    "image": "dwhitena/shelfrunner",
    "payload": "{\"asset\": \"579f803559b0b30009ff60f1\", \"mongo_host\": \"146.148.41.180\"}"
  }]
}' "http://146.148.41.180:8080/v1/groups/mygroup/jobs"
```
