Example output:

```
NUMBER OF DOCUMENTS: 100000
NUMBER OF REQUESTS PER QUERY: 100

=============================================
QUERY #1: All comments authored by a random user
NUMBER OF RELATIONSHIP LEVELS: 1

MONGO MEAN RESPONSE TIME (seconds): 0.0886
GRAPH MEAN RESPONSE TIME (seconds): 0.3520

=============================================
QUERY #2: All assets commented on by a user
NUMBER OF RELATIONSHIP LEVELS: 1

MONGO MEAN RESPONSE TIME (seconds): 0.0891
GRAPH MEAN RESPONSE TIME (seconds): 0.3677

=============================================
QUERY #3: All comments on an asset along with corresponding authors
NUMBER OF RELATIONSHIP LEVELS: 1-2

MONGO MEAN RESPONSE TIME (seconds): 0.0897
GRAPH MEAN RESPONSE TIME (seconds): 0.3745

=============================================
QUERY #4: All comments parented by a set of comments, the set of comments being authored by the autho
r of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 4

MONGO MEAN RESPONSE TIME (seconds): 0.1864
GRAPH MEAN RESPONSE TIME (seconds): 0.3618

=============================================
QUERY #5: All comments grandparented by a set of comments, the set of comments being authored by the 
author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 5

MONGO MEAN RESPONSE TIME (seconds): 0.2811
GRAPH MEAN RESPONSE TIME (seconds): 0.3684

=============================================
QUERY #6: All comments great-grandparented by a set of comments, the set of comments being authored b
y the author of the parent of a given comment
NUMBER OF RELATIONSHIP LEVELS: 6

MONGO MEAN RESPONSE TIME (seconds): 0.3761
GRAPH MEAN RESPONSE TIME (seconds): 0.3631
```
