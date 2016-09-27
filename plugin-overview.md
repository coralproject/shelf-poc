# plugin-poc

Proof of concept home for Coral's plugin architecture

## Goals
To provide an approachable way to build on the Coral Platform:

* Add pages to Cay
* Create embeddable widgets
* Install / customize / swap out front end components
* Define Item Types (exposing corresponding endpoints)
* Define Item Relationships 
* Define Views (and kick off their creation)
* Define Endpoints
  * Queries (aggregation pipelines on collections or views)
  * Algorithms (written in any language, served by iron.io)

  
## Installer

Installing a component:

### Back End

1. Register Item Types (Sponge)
  2. Initialize CRUD endpoints for the Item Type
1. Register Item Relationships (Shelf)
  1. Scan existing Items to populate relationships
1. Register Views (Shelf)
  2. Install graph view configuration
  3. Install upload worker views through Shelf into iron.io
  4. Kick off processes that create Views marked for pre-rendering  
1. Register Xenia Query Sets
  1. Aggregation Pipeline endpoints
  2. Worker endpoints 
  3. Transpiler based enpoints? (GraphQL)
 
### Front End

1. Define import paths for plugin components by item type
  2.  ie, ```import Coral.ItemType('comment'); ``` resolves to importing the comment from this plugin
3. Add Routes to install pages into Cay
  4. Kick off a build incorporating this plugin 
5. Kick off builds for embeddable components


 
