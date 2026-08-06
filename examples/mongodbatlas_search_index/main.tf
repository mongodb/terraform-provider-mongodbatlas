# Text search index with dynamic mappings.
resource "mongodbatlas_search_index" "search" {
  project_id       = var.project_id
  cluster_name     = var.cluster_name
  database         = var.database
  collection_name  = var.collection_name
  name             = "search-index"
  type             = "search"
  mappings_dynamic = true
}

# Vector search index with embeddings you generate and store yourself.
resource "mongodbatlas_search_index" "vector" {
  project_id      = var.project_id
  cluster_name    = var.cluster_name
  database        = var.database
  collection_name = var.collection_name
  name            = "vector-index"
  type            = "vectorSearch"
  fields          = <<-EOF
[{
  "type": "vector",
  "path": "plot_embedding",
  "numDimensions": 1536,
  "similarity": "euclidean"
}]
EOF
}

# Vector search index with Automated Embedding: Atlas generates the embeddings
# from a text field using the specified Voyage model.
resource "mongodbatlas_search_index" "auto_embed" {
  project_id      = var.project_id
  cluster_name    = var.cluster_name
  database        = var.database
  collection_name = var.collection_name
  name            = "auto-embed-index"
  type            = "vectorSearch"
  fields          = <<-EOF
[{
  "type": "autoEmbed",
  "path": "plot",
  "model": "voyage-4-lite",
  "modality": "text"
},
{
  "type": "filter",
  "path": "genres"
}]
EOF
}

output "search_index_id" {
  value = mongodbatlas_search_index.search.index_id
}

output "vector_index_id" {
  value = mongodbatlas_search_index.vector.index_id
}

output "auto_embed_index_id" {
  value = mongodbatlas_search_index.auto_embed.index_id
}
