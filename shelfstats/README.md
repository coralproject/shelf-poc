## Example output for 100k documents

```
NUMBER OF DOCUMENTS: 100000
NUMBER OF REQUESTS PER QUERY: 10

=============================================
QUERY #1: All comments authored by a random user
NUMBER OF RELATIONSHIP LEVELS: 1

MONGO MEAN RESPONSE TIME (seconds): 0.1604
GRAPH MEAN RESPONSE TIME (seconds): 0.4108

=============================================
QUERY #2: All assets commented on by a user
NUMBER OF RELATIONSHIP LEVELS: 1

MONGO MEAN RESPONSE TIME (seconds): 0.1620
GRAPH MEAN RESPONSE TIME (seconds): 0.4326

=============================================
QUERY #3: All comments on an asset along with corresponding authors
NUMBER OF RELATIONSHIP LEVELS: 1-2

MONGO MEAN RESPONSE TIME (seconds): 0.1689
GRAPH MEAN RESPONSE TIME (seconds): 0.4356

=============================================
QUERY #4: All comments parented by a set of comments, the set of comments being authored by the author 
of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 4

MONGO MEAN RESPONSE TIME (seconds): 0.2674
GRAPH MEAN RESPONSE TIME (seconds): 0.4206

=============================================
QUERY #5: All comments grandparented by a set of comments, the set of comments being authored by the 
author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 5

MONGO MEAN RESPONSE TIME (seconds): 0.3458
GRAPH MEAN RESPONSE TIME (seconds): 0.4683

=============================================
QUERY #6: All comments great-grandparented by a set of comments, the set of comments being authored by 
the author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 6

MONGO MEAN RESPONSE TIME (seconds): 0.4676
GRAPH MEAN RESPONSE TIME (seconds): 0.4532
```

## Example output for 1M documents

Note, in this case, documents were directly uploaded without embedded Mongo relationships.  Thus, the Mongo comparison numbers are left out.

```
NUMBER OF DOCUMENTS: 1000000
NUMBER OF REQUESTS PER QUERY: 10

=============================================
QUERY #1: All comments authored by a random user
NUMBER OF RELATIONSHIP LEVELS: 1

GRAPH MEAN RESPONSE TIME (seconds): 3.7784

=============================================
QUERY #2: All assets commented on by a user
NUMBER OF RELATIONSHIP LEVELS: 1

GRAPH MEAN RESPONSE TIME (seconds): 3.7782

=============================================
QUERY #3: All comments on an asset along with corresponding authors
NUMBER OF RELATIONSHIP LEVELS: 1-2

GRAPH MEAN RESPONSE TIME (seconds): 3.7915

=============================================
QUERY #4: All comments parented by a set of comments, the set of comments being authored by the author 
of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 4

GRAPH MEAN RESPONSE TIME (seconds): 3.7482

=============================================
QUERY #5: All comments grandparented by a set of comments, the set of comments being authored by the 
author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 5

GRAPH MEAN RESPONSE TIME (seconds): 3.7074

=============================================
QUERY #6: All comments great-grandparented by a set of comments, the set of comments being authored by 
the author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 6

GRAPH MEAN RESPONSE TIME (seconds): 3.6712
```

## Example output for 10M documents

Note, in this case, documents were directly uploaded without embedded Mongo relationships.  Thus, the Mongo comparison numbers are left out.

```
NUMBER OF DOCUMENTS: 10000000
NUMBER OF REQUESTS PER QUERY: 10

=============================================
QUERY #1: All comments authored by a random user
NUMBER OF RELATIONSHIP LEVELS: 1

GRAPH MEAN RESPONSE TIME (seconds): 32.6445

=============================================
QUERY #2: All assets commented on by a user
NUMBER OF RELATIONSHIP LEVELS: 1

GRAPH MEAN RESPONSE TIME (seconds): 34.7760

=============================================
QUERY #3: All comments on an asset along with corresponding authors
NUMBER OF RELATIONSHIP LEVELS: 1-2

GRAPH MEAN RESPONSE TIME (seconds): 34.8924

=============================================
QUERY #4: All comments parented by a set of comments, the set of comments being authored by the author 
of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 4

GRAPH MEAN RESPONSE TIME (seconds): 33.8968

=============================================
QUERY #5: All comments grandparented by a set of comments, the set of comments being authored by the 
author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 5

GRAPH MEAN RESPONSE TIME (seconds): 33.7320

=============================================
QUERY #6: All comments great-grandparented by a set of comments, the set of comments being authored 
by the author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 6

GRAPH MEAN RESPONSE TIME (seconds): 33.8621
```
