# POC for "Shelf"

The goals of shelf are:

1. To enable utilization a single collection of items for all assets, authors, comments, actions, etc., such that management of the items is straightforward and such that we don't need to worry about complex indexing of the items for joins.
2. To generate on-the-fly or cached sub-collections of items (or "views") for analysis, aggregation, display, etc.
3. To manage all relationships between items ("authored by," "blocked," etc.).

This POC mocks out much of this functionality, such that we can evaluate architecture decisions, formatting, config, etc.  The POC includes the following:

- Data importing/ingestion using a branch of `xenia` ([item-cayley](https://github.com/coralproject/xenia/tree/item-cayley)) with a modified internal `item` package and `sponged` cmd.
- Dummy data and relationship generation using the `dummy` program included [here](dummy).
- Querying based on relationships using the `relquery` web server included [here](relquery).

**Disclaimer** - The code in this repo is meant for evaluation only and should not be utilized in any production systems.  After evaluating functionality with this POC, actual shelf functionality will be implemented as part of the xenia repo.

## Setting up your environment, usage

### Dependencies:

  1. `sponged` - as built from the `item-cayley` branch of `xenia`
  2. MongoDB - note, you may want a separate instance of mongo running for this POC as it will create databases, indices, etc.  Just so it doesn't break anything you are currently using that is dependent on Mongo.

### 1. Prepare data

If you want to create dummy data for testing:

  1. Build the `dummy` binary by executing `go build` [here](dummy), OR build a Docker image including the binary using the [Makefile](dummy/Makefile).  Note, you would likely want to change the tags of the docker image to build locally or push up to the docker registry of your choice.
  2. Copy the [dummy.env](dummy/dummy.env) file to `/etc` and modify it according to where `sponged` is installed and how many documents you want to generate.
  3. Run `dummy` (or the Docker image if you built it).  When this executes, it will:
    - Generate dummy items including assets, comments, and users in the proportions 5%, 80%, and 15%, respectively.
    - For each of the generated items, send the item to `sponged`, where `sponged` will format, verify, and store the item along with inserting the item and any relevant relationships into the Cayley graph DB.

If you want to utilize real world data for testing:

  1. Make sure `sponged` (based on the `item-cayley` branch of xenia) is running.
  2. Create a script that push assets, comments, and authors to the `sponged` endpoint defined [here](https://github.com/coralproject/xenia/blob/item-cayley/cmd/sponged/routes/routes.go#L108).
  3. Run your script to import the data.

### 2. Query the dummy data

1. Build the `relquery` binary by executing `go build` [here](relquery), OR build a Docker image including the binary using the [Makefile](relquery/Makefile).  Again, you would likely want to change the tags of the docker image to build locally or push up to the docker registry of your choice.
2. Run `relquery`.

You will now be able to make use of the following endpoints:

- GET `/asset?num=<number of assets>` - this endpoint will retrieve a random asset item from MongoDB, given the total number of assets from which a random sample will be pulled.
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

- GET `/graph?asset=<asset id>` - this endpoint will use the graph DB to retrieve all comments related to the given asset along with the users authoring the comments.
  - Example request: 
    
    ```
    GET http://<host>:8080/graph?asset=578f5324b1df410001e01da0
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

- GET `/mongo?asset=<asset id>` - this endpoint will use only mongo to retrieve all comments related to the given asset along with the users authoring the comments.
  - Example request: 
    
    ```
    GET http://<host>:8080/mongo?asset=578f5324b1df410001e01da0
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
