# POC for "Shelf"

The goals of shelf are:

1. To enable utilization a single collection of items for all assets, authors, comments, actions, etc., such that management of the items is straightforward and such that we don't need to worry about complex indexing of the items for joins.
2. To generate on-the-fly or cached sub-collections of items (or "views") for analysis, aggregation, display, etc.
3. To manage all relationships between items ("authored by," "blocked," etc.).

This POC mocks out much of this functionality, such that we can evaluate architecture decisions, formatting, config, etc.  The POC includes the following:

- Data importing/ingestion using a branch of `xenia` ([item-cayley](https://github.com/coralproject/xenia/tree/item-cayley)) with a modified internal `item` package and `sponged` cmd.
- Dummy data and relationship generation using the `shelfdummy-direct` program included [here](shelfdummy-direct) or the `shelfdummy-sponged` included [here](shelfdummy-sponged).
- View generation based on relationships using the `shelfquery` web server included [here](shelfquery).
- Stats collection for views using the `shelfstats` script includes [here](shelfstats).

**Disclaimer** - The code in this repo is meant for evaluation only and should not be utilized in any production systems.  After evaluating functionality with this POC, actual shelf functionality will be implemented as part of the xenia repo.

## Setting up your environment, usage

### Dependencies:

  1. `sponged` - as built from the `item-cayley` branch of `xenia`
  2. MongoDB - note, you may want a separate instance of mongo running for this POC as it will create databases, indices, etc.  Just so it doesn't break anything you are currently using that is dependent on Mongo.

### 1. Prepare data

If you want to create dummy data for testing using directly, where items do not have embedded relationships and relationships are only reflected in the graph DB:

  1. Build the `shelfdummy-direct` binary by executing `go build` [here](shelfdummy-direct), OR build a Docker image including the binary using the [Makefile](shelfdummy-direct/Makefile).  Note, you would likely want to change the tags of the docker image to build locally or push up to the docker registry of your choice.
  2. Define the example environmental variables in [shelf.env](files/shelf.env).  Modify them according to how many documents you want to generate.
  3. Run `shelfdummy-direct` (or the Docker image if you built it).  When this executes, it will:
    - Generate dummy items including assets, comments, and users in the proportions 5%, 80%, and 15%, respectively.
    - For each of the generated items, store the item along with inserting the item and any relevant relationships into the Cayley graph DB.

If you want to create dummy data for testing using `sponged`, where `sponged` will both embed relationships in Mongo and create graph DB items:

  1. Build the `shelfdummy-sponged` binary by executing `go build` [here](shelfdummy-sponged), OR build a Docker image including the binary using the [Makefile](shelfdummy-sponged/Makefile).  Note, you would likely want to change the tags of the docker image to build locally or push up to the docker registry of your choice.
  2. Define the example environmental variables in [dummy.env](files/dummy.env).  Modify them according to where `sponged` is installed and how many documents you want to generate.
  3. Run `shelfdummy-sponged` (or the Docker image if you built it).  When this executes, it will:
    - Generate dummy items including assets, comments, and users in the proportions 5%, 80%, and 15%, respectively.
    - For each of the generated items, send the item to `sponged`, where `sponged` will format, verify, and store the item along with inserting the item and any relevant relationships into the Cayley graph DB.

If you want to utilize real world data for testing:

  1. Make sure `sponged` (based on the `item-cayley` branch of xenia) is running.
  2. Create a script that push assets, comments, and authors to the `sponged` endpoint defined [here](https://github.com/coralproject/xenia/blob/item-cayley/cmd/sponged/routes/routes.go#L108).
  3. Run your script to import the data.

### 2. Generate views (query the dummy data)

1. Build the `shelfquery` binary by executing `go build` [here](shelfquery), OR build a Docker image including the binary using the [Makefile](shelfquery/Makefile).  Again, you would likely want to change the tags of the docker image to build locally or push up to the docker registry of your choice.
2. Run `shelfquery`.

You will now be able to make use of the view endpoints described [here](shelfquery/README.md).

### 3. Gather stats

1. Build the `shelfstats` binary by executing `go build` [here](shelfstats).
2. Set [these](files/stats.env) environmental vars.
3. Run `shelfstats`.  This will output stats (see example [here](shelfstats/README.md)) for all available queries to standard out.
