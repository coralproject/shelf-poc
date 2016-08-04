# Performance and Functionality Analysis

Sponge, Shelf and Xenia make up the data handling layer. The following benchmarks are designed to simulate the core use cases and look forward to reach goals.

Once an end to end prototype is completed, these use cases can be tested to predict the application's behavior under a variety of simulated real world situations.

In addition to applicatoin benchmarks, this page includes a series of features representative of the reach goals that we hope this architecture will allow.

## Benchmarks

### Comment threads on assets

#### Input

* Comments/s (comments inputted per second)
* Actions/s (actions randomly taken on new comments)
* Assets (number of assets the comments are distributed over)


#### Views

* "Comment Thread on Asset"
  * Asset (id)
  * Comments "on" the asset
  * Users who "authored" comments "on" the asset
  * Actions "on" the comments on the asset
  * Users performing Actions "on" the comments on the asset
  * (Users "mentioned" in comments on assets)

  
#### Queries

* Gets (all top level assets and counts for their children)
* My activity (all replies and likes to comments I have written)
* Count (number of comments and likes and unique users in thread)

#### Measurements

* Read and write response times
* How long does a new comment or action take to make it from input to output?
* How does the system perform as Item Store size varies?

## Personalization Modifyers 

We will want to apply a series of modifyers on top of the basic benchmark behavior above so that we can deliver a uniquely cusomized experience for each user. 

Ultimately, enabling this level of insight and personalization is the justification for designing the data layer at the item/relationship/view abstraction as opposed to building straightforward optimized apis.

### Structure

Modifyers consist of 2 parts, 

```(graph search) -> (concept)```

The graph query describes the items we're talking about.

The applied concept says what we do with those items.

Depending on many factors, we may want to treat these searches as Views or perform them as needed.

### Examples

These are _reach goals_, not hard requirements at scale, but representative of the kinds of things that we are trying to achieve with the platform.


* If you have someone you have muted who likes someone's comments who have also liked theirs back -> you can mute them.
* Notification mentions will only come up if you follow the user who mentions you
* If you have likes more than x percent of someone's comments on y topic (topics trickle down from assets) in z time, then the system will recommend their next comments on any asset on that topic.
* If you have had a reply from the Author of an asset and you liked that reply, you get a notification when that author's next article is published.
* Notify if someone replies to one of your comments if the user who wrote the comment is followed by someone you follow, or you directly follow and are not blocking.  (up to x degrees of Kevin Bacon.)
* If you @ someone who follows you, start a private conversation within the public comment space (possible unanimous vote to make it public?)
* Comments -> Commentary (tag own comments and build a page showing them in context of articles on a timeline)

### Meta

Because simple crazy ideas aren't good enough

* Compose the above together to form meta scores (Trust for Graph relationships)
* Entropy algorithms based on your interactions with a user (activity, not time)!!


### Other product features discussed

* Mechanism to store individual comments people have seen. (rolling stats digests/delete?)
  * Tab: comments I haven't seen on this thread (that we think you want to see)
* Collaboration on comments



