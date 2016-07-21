## GET `/asset?num=<number of assets>` 

This endpoint will retrieve a random asset item from MongoDB, given the total number of assets from which a random sample will be pulled.
  - Example request: 
  
    ```
    GET http://<host>:8080/asset?num=500
    ```
  
  - Example response:
  
    ```
    {
      "id": "578f5324b1df410001e01d9a",
      "t": "coral_asset",
      "v": 1,
      "d": {
        "id": 248,
        "name": "lgXlnHWIRXxeJwzYTaYO"
      }
    }
    ```

## GET `/comment?num=<number of comments>` 

This endpoint will retrieve a random comment item from MongoDB, given the total number of comments from which a random sample will be pulled.
  - Example request: 
  
    ```
    GET http://<host>:8080/comment?num=500
    ```
  
  - Example response:
  
    ```
    {
      "id": "578fd169b1df410001d43190",
      "t": "coral_comment",
      "v": 1,
      "d": {
        "asset_id": 3640,
        "body": "YjjVmUlsdjsuMaRbEpCr",
        "id": 12736,
        "user_id": 13340
      },
      "rels": [
        {
          "n": "context",
          "t": "coral_asset",
          "id": "578fd03bb1df410001d3fa80"
        },
        {
          "n": "author",
          "t": "coral_user",
          "id": "578fd024b1df410001d3e5cc"
        }
      ]
    }
    ```

## GET `/user?num=<number of users>` 

This endpoint will retrieve a random user item from MongoDB, given the total number of users from which a random sample will be pulled.
  - Example request: 
  
    ```
    GET http://<host>:8080/user?num=500
    ```
  
  - Example response:
  
    ```
    {
      "id": "578fd006b1df410001d3b81f",
      "t": "coral_user",
      "v": 1,
      "d": {
        "id": 1647,
        "name": "KsIxvlZayj"
      }
    }    
    ```

## GET `/graph/singleasset?asset=<asset id>` 

This endpoint will use the graph DB to retrieve all comments related to the given asset along with the users authoring the comments.
  - Example request: 
    
    ```
    GET http://<host>:8080/graph/singleasset?asset=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f5321b1df410001e01716",
        "t": "coral_user",
        "v": 1,
        "d": {
          "id": 86,
          "name": "hskmHIyABE"
        }
      },
      {
        "id": "578f5321b1df410001e01758",
        "t": "coral_user",
        "v": 1,
        "d": {
          "id": 147,
          "name": "AIFXsRBOdv"
        }
      },
      {
        "id": "578f5339b1df410001e02945",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "nTFcpRoUEzrMPWSRWjlf",
          "id": 2740,
          "user_id": 1382
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c2c"
          }
        ]
      },
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```

## GET `/mongo/singleasset?asset=<asset id>` 

This endpoint will use only mongo to retrieve all comments related to the given asset along with the users authoring the comments.
  - Example request: 
    
    ```
    GET http://<host>:8080/mongo/singleasset?asset=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f5321b1df410001e01716",
        "t": "coral_user",
        "v": 1,
        "d": {
          "id": 86,
          "name": "hskmHIyABE"
        }
      },
      {
        "id": "578f5321b1df410001e01758",
        "t": "coral_user",
        "v": 1,
        "d": {
          "id": 147,
          "name": "AIFXsRBOdv"
        }
      },
      {
        "id": "578f5339b1df410001e02945",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "nTFcpRoUEzrMPWSRWjlf",
          "id": 2740,
          "user_id": 1382
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c2c"
          }
        ]
      },
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```

## GET `/graph/userassets?user=<user id>` 

This endpoint will use the graph DB to retrieve all assets commented on by a given user.
  - Example request: 
    
    ```
    GET http://<host>:8080/graph/userassets?user=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578fd035b1df410001d3f5e5",
        "t": "coral_asset",
        "v": 1,
        "d": {
          "id": 2461,
          "name": "iQOCUfCSKWQeDiqrIjmv"
        }
      },
      {
        "id": "578fd039b1df410001d3f905",
        "t": "coral_asset",
        "v": 1,
        "d": {
          "id": 3261,
          "name": "sXCbqwZYeiBltmJyVYpA"
        }
      }
    ]
    ```

## GET `/mongo/userassets?user=<user id>` 

This endpoint will use mongo DB only to retrieve all assets commented on by a given user.
  - Example request: 
    
    ```
    GET http://<host>:8080/mongo/userassets?user=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578fd035b1df410001d3f5e5",
        "t": "coral_asset",
        "v": 1,
        "d": {
          "id": 2461,
          "name": "iQOCUfCSKWQeDiqrIjmv"
        }
      },
      {
        "id": "578fd039b1df410001d3f905",
        "t": "coral_asset",
        "v": 1,
        "d": {
          "id": 3261,
          "name": "sXCbqwZYeiBltmJyVYpA"
        }
      }
    ]
    ```

## GET `/graph/usercomments?user=<user id>` 

This endpoint will use the graph DB to retrieve all comments authored by the given user.
  - Example request: 
    
    ```
    GET http://<host>:8080/graph/usercomments?user=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```

## GET `/mongo/usercomments?user=<user id>` 

This endpoint will use mongo only to retrieve all comments authored by the given user.
  - Example request: 
    
    ```
    GET http://<host>:8080/mongo/usercomments?user=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```

## GET `/graph/parentcomments?comment=<comments id>` 

This endpoint will use the graph DB to retrieve all comments parented by comments authored by the author of the parent of the provided comment.
  - Example request: 
    
    ```
    GET http://<host>:8080/graph/parentcomments?comment=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```

## GET `/mongo/parentcomments?comment=<comments id>` 

This endpoint will use mongo only to retrieve all comments parented by comments authored by the author of the parent of the provided comment.
  - Example request: 
    
    ```
    GET http://<host>:8080/mongo/parentcomments?comment=578f5324b1df410001e01da0
    ```
  
  - Example response:
  
    ```
    [
      {
        "id": "578f533eb1df410001e02b73",
        "t": "coral_comment",
        "v": 1,
        "d": {
          "asset_id": 258,
          "body": "rEdkZHdlsGOdVpkGdZvJ",
          "id": 3296,
          "parent_id": 2003,
          "user_id": 1458
        },
        "rels": [
          {
            "n": "context",
            "t": "coral_asset",
            "id": "578f5324b1df410001e01da0"
          },
          {
            "n": "author",
            "t": "coral_user",
            "id": "578f5324b1df410001e01c73"
          },
          {
            "n": "parent",
            "t": "coral_comment",
            "id": "578f5333b1df410001e02665"
          }
        ]
      }
    ]
    ```


