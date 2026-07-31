# MongoDB Atlas Provider - Search Indexes on a Cluster

This example shows how to create the different `mongodbatlas_search_index` types on an Atlas cluster. A project and cluster are created as a prerequisite, and then three indexes on the same collection:

- A text `search` index with dynamic mappings.
- A `vector` index, for embeddings you generate and store yourself.
- An `autoEmbed` (Automated Embedding) index, where Atlas generates the embeddings from a text field using a Voyage model.

Variables required to be set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret.
- `org_id`: Organization ID where the project and cluster are created.
- `database`: Database that holds the collection to index.
- `collection_name`: Collection to index.

Adjust the field paths in `main.tf` to match your collection.

For the Automated Embedding (`autoEmbed`) index, see [Automated Embedding](https://www.mongodb.com/docs/vector-search/crud-embeddings/automated-embedding/) for supported models and requirements. For additional information on search index field types, see [Vector Search Types](https://www.mongodb.com/docs/atlas/atlas-vector-search/vector-search-type/).
