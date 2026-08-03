# MongoDB Atlas Provider - Search Indexes on a Cluster

This example shows how to create different `mongodbatlas_search_index` types on a MongoDB Atlas cluster, three indexes on the same collection:

- A text `search` index with dynamic mappings.
- A `vector` index, for embeddings you generate and store yourself.
- An `autoEmbed` (Automated Embedding) index, where Atlas generates the embeddings from a text field using a Voyage model.

Variables required to be set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret.
- `project_id`: MongoDB Atlas project ID that holds the cluster.
- `cluster_name`: MongoDB Atlas cluster with the collection to index.
- `database`: Database that holds the collection to index.
- `collection_name`: Collection to index.

Adjust the field paths in `main.tf` to match your collection.

For the Automated Embedding (`autoEmbed`) index, see [Automated Embedding](https://www.mongodb.com/docs/vector-search/crud-embeddings/automated-embedding/) for supported models and requirements. For additional information on search index field types, see [Vector Search Types](https://www.mongodb.com/docs/atlas/atlas-vector-search/vector-search-type/).
