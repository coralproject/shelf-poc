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
