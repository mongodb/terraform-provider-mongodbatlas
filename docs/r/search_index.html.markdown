# Resource: mongodbatlas_search_index

`mongodbatlas_search_index` provides a Search Index resource. This allows indexes to be created.

## Example Usage

### Basic search index
```terraform
resource "mongodbatlas_search_index" "test-basic-search-index" {
  name   = "test-basic-search-index"
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  
  analyzer = "lucene.standard"
  collection_name = "collection_test"
  database = "database_test"
  mappings_dynamic = true
  
  search_analyzer = "lucene.standard"
}
```

### Basic vector index
```terraform
resource "mongodbatlas_search_index" "test-basic-search-vector" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  collection_name = "collection_test"
  database = "database_test"
  type = "vectorSearch"
  fields = <<-EOF
[{
      "type": "vector",
      "path": "plot_embedding",
      "numDimensions": 1536,
      "similarity": "euclidean"
}]
EOF
}
```

### Advanced search index (with custom analyzers)
```terraform
resource "mongodbatlas_search_index" "test-advanced-search-index" {
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  analyzer = "lucene.standard"
  collection_name = "collection_test"
  database = "database_test"
  mappings_dynamic = false
  mappings_fields = <<-EOF
{
      "address": {
        "type": "document",
        "fields": {
          "city": {
            "type": "string",
            "analyzer": "lucene.simple",
            "ignoreAbove": 255
          },
          "state": {
            "type": "string",
            "analyzer": "lucene.english"
          }
        }
      },
      "company": {
        "type": "string",
        "analyzer": "lucene.whitespace",
        "multi": {
          "mySecondaryAnalyzer": {
            "type": "string",
            "analyzer": "lucene.french"
          }
        }
      },
      "employees": {
        "type": "string",
        "analyzer": "lucene.standard"
      }
}
EOF
  name = "test-advanced-search-index"
  search_analyzer = "lucene.standard"
  analyzers = <<-EOF
  [{
  "name": "index_analyzer_test_name",
  "charFilters": [{
	"type": "mapping",
	"mappings": {"\\" : "/"}
    	}],
  "tokenizer": {
  "type": "nGram",
  "minGram": 2,
  "maxGram": 5
		},
  "tokenFilters": [{
	"type": "length",
	"min": 20,
	"max": 33
    	}]
  }]
EOF
  synonyms {
    analyzer = "lucene.simple"
    name = "synonym_test"
    source_collection = "collection_test"
  }
}
```

## Argument Reference

* `type` - (Optional) Type of index: `search` or `vectorSearch`. Default type is `search`.
* `name` - (Required) The name of the search index you want to create.
* `project_id` - (Required) The ID of the organization or project you want to create the search index within.
* `cluster_name` - (Required) The name of the cluster where you want to create the search index within.
* `wait_for_index_build_completion` - (Optional) Wait for search index to achieve Active status before terraform considers resource built.
* `timeouts`- (Optional) The duration of time to wait for Search Index to be created, updated, or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Serach Index create & update is `3h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).


* `analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when creating the index. Defaults to [lucene.standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/standard/#std-label-ref-standard-analyzer)

* `analyzers` - [Custom analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) to use in this index. This is an array of JSON objects.
```
analyzers = <<-EOF
  [{
  "name": "index_analyzer_test_name",
  "charFilters": [{
	"type": "mapping",
	"mappings": {"\\" : "/"}
    	}],
  "tokenizer": {
  "type": "nGram",
  "minGram": 2,
  "maxGram": 5
	},
  "tokenFilters": [{
	"type": "length",
	"min": 20,
	"max": 33
    	}]
  }]
EOF
```

* `collection_name` - (Required) Name of the collection the index is on.

* `database` - (Required) Name of the database the collection is in.

* `mappings_dynamic` - Indicates whether the search index uses dynamic or static mapping. For dynamic mapping, set the value to `true`. For static mapping, specify the fields to index using `mappings_fields`

* `mappings_fields` - attribute is required in search indexes when `mappings_dynamic` is false. This field needs to be a JSON string in order to be decoded correctly.
  ```terraform
    mappings_fields = <<-EOF
    {
      "address": {
        "type": "document",
        "fields": {
          "city": {
            "type": "string",
            "analyzer": "lucene.simple",
            "ignoreAbove": 255
          },
          "state": {
            "type": "string",
            "analyzer": "lucene.english"
          }
        }
      },
      "company": {
        "type": "string",
        "analyzer": "lucene.whitespace",
        "multi": {
          "mySecondaryAnalyzer": {
            "type": "string",
            "analyzer": "lucene.french"
          }
        }
      },
      "employees": {
        "type": "string",
        "analyzer": "lucene.standard"
      }
    }
    EOF
  ```

* `search_analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when searching the index. Defaults to [lucene.standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/standard/#std-label-ref-standard-analyzer)
* `synonyms` - Synonyms mapping definition to use in this index.

* `fields` - Array of [Fields](https://www.mongodb.com/docs/atlas/atlas-search/field-types/knn-vector/#std-label-fts-data-types-knn-vector) to configure this `vectorSearch` index. It is mandatory for vector searches and it must contain at least one `vector` type field. This field needs to be a JSON string in order to be decoded correctly.

* `stored_source` - String that can be "true" (store all fields), "false" (default, don't store any field), or a JSON string that contains the list of fields to store (include) or not store (exclude) on Atlas Search. To learn more, see [Stored Source Fields](https://www.mongodb.com/docs/atlas/atlas-search/stored-source-definition/).
  ```terraform
    stored_source = <<-EOF
    {
      "include": ["field1", "field2"]
    }
    EOF
  ```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `index_id` - The unique identifier of the Atlas Search index.
* `status` - Current status of the index.

### Analyzers (search  index)
An [Atlas Search analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/) prepares a set of documents to be indexed by performing a series of operations to transform, filter, and group sequences of characters. You can define a custom analyzer to suit your specific indexing needs.

* `name` - (Required) 	
  Name of the custom analyzer. Names must be unique within an index, and may **not** start with any of the following strings:
    * `lucene`
    * `builtin`
    * `mongodb`
* `charFilters` - Array containing zero or more character filters. Always require a `type` field, and some take additional options as well
  ```terraform
  "charFilters":[{
   "type": "<FILTER_TYPE>",
   "ADDITIONAL_OPTION": VALUE
  }]
  ```
  Atlas search supports four `types` of character filters:
    * [htmlStrip](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-htmlStrip-ref) - Strips out HTML constructs
        * `type` - (Required) Must be `htmlStrip`
        * `ignoredTags`- a list of HTML tags to exclude from filtering
        ```terraform
          analyzers = <<-EOF [{
            "name": "analyzer_test",
            "charFilters":[{
              "type": "htmlStrip",
              "ignoredTags": ["a"]
              }]   
            }] 
       ```
    * [icuNormalize](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-icuNormalize-ref) - Normalizes text with the [ICU](http://site.icu-project.org/) Normalizer. Based on Lucene's [ICUNormalizer2CharFilter](https://lucene.apache.org/core/8_3_0/analyzers-icu/org/apache/lucene/analysis/icu/ICUNormalizer2CharFilter.html)
    * [mapping](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-mapping-ref) - Applies user-specified normalization mappings to characters. Based on Lucene's [MappingCharFilter](https://lucene.apache.org/core/8_0_0/analyzers-common/org/apache/lucene/analysis/charfilter/MappingCharFilter.html)
      An object containing a comma-separated list of mappings. A mapping indicates that one character or group of characters should be substituted for another, in the format `<original> : <replacement>`
      ### Example
        ```terraform
        analyzers = <<-EOF [{
          "name":"name_analyzer",        
          "type": "mapping",
          "mappings":  
          {
             "\\" : "/"
          }
          }]
          EOF 
      ```
    * [persian](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-persian-ref) - Replaces instances of [zero-width non-joiners](https://en.wikipedia.org/wiki/Zero-width_non-joiner) with ordinary space. Based on Lucene's [PersianCharFilter](https://lucene.apache.org/core/8_0_0/analyzers-common/org/apache/lucene/analysis/fa/PersianCharFilter.html)


* `tokenizer` - (Required) Tokenizer to use in search indexes. Determines how Atlas Search splits up text into discrete chunks of indexing. Always require a type field, and some take additional options as well.
    ```terraform
    "tokenizer":{
    "type": "<tokenizer-type>",
    "ADDITIONAL_OPTIONS": VALUE
    }
    ```
  Atlas Search supports the following tokenizer options:
    * [standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-standard-tokenizer-ref) - Tokenize based on word break rules from the [Unicode Text Segmentation algorithm](http://www.unicode.org/L2/L2019/19034-uax29-34-draft.pdf):
        * `type` - Must be `standard`
        * `maxTokenLength` - Maximum length for a single token. Tokens greater than this length are split at `maxTokenLength` into multiple tokens.
    * [keyword](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-keyword-tokenizer-ref) - Tokenize the entire input as a single token.
    * [whitespace](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-whitespace-tokenizer-ref) - Tokenize based on occurrences of whitespace between words.
        * `type` - Must be `whitespace`
        * `maxTokenLength` - Maximum length for a single token. Tokens greater than this length are split at `maxTokenLength` into multiple tokens.
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tokenizer-ref) - Tokenize into text chunks, or "n-grams", of given sizes.
        * `type` - Must be `nGram`
        * `minGram` - (Required) Number of characters to include in the shortest token created.
        * `maxGram` - (Required) Number of characters to include in the longest token created.
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tokenizer-ref) - Tokenize input from the beginning, or "edge", of a text input into n-grams of given sizes.
        * `type` - Must be `edgeGram`
        * `minGram` - (Required) Number of characters to include in the shortest token created.
        * `maxGram` - (Required) Number of characters to include in the longest token created.
    * [regexCaptureGroup](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-regexcapturegroup-tokenizer-ref) - Match a regular expression pattern to extract tokens.
        * `type` - Must be `regexCaptureGroup`
        * `pattern` - (Required) A regular expression to match against.
        * `group` - (Required) Index of the character group within the matching expression to extract into tokens. Use 0 to extract all character groups.
    * [regexSplit](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-regexSplit-tokenizer-ref) - Split tokens with a regular-expression based delimiter.
        * `type` - Must be `regexSplit`
        * `pattern` - (Required) A regular expression to match against.
    * [uaxUrlEmail](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-uaxUrlEmail-tokenizer-ref) - Tokenize URLs and email addresses. Although [uaxUrlEmail](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-uaxUrlEmail-tokenizer-ref) tokenizer tokenizes based on word break rules from the [Unicode Text Segmentation algorithm](http://www.unicode.org/L2/L2019/19034-uax29-34-draft.pdf), we recommend using uaxUrlEmail tokenizer only when the indexed field value includes URLs and email addresses. For fields that do not include URLs or email addresses, use the [standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-standard-tokenizer-ref) tokenizer to create tokens based on word break rules.
        * `type` - Must be `uaxUrlEmail`
        *  `maxTokenLength` - The maximum number of characters in one token.

* `token_filters` - Array containing zero or more token filters. Always require a type field, and some take additional options as well:
  ```terraform
  "tokenFilters":[{
    "type": "<FILTER_TYPE>",
    "ADDITIONAL-OPTIONS": VALUE
  }]
  ```
  Atlas Search supports the following token filters:
    * [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) - Creates tokens for words that sound the same based on [Daitch-Mokotoff Soundex](https://en.wikipedia.org/wiki/Daitch%E2%80%93Mokotoff_Soundex) phonetic algorithm. This filter can generate multiple encodings for each input, where each encoded token is a 6 digit number:
        *  `type` - Must be `daitchMokotoffSoundex`
        * `originalTokens` - Specifies whether to include or omit the original tokens in the output of the token filter. Value can be one of the following:
            * `include` - to include the original tokens with the encoded tokens in the output of the token filter. We recommend this value if you want queries on both the original tokens as well as the encoded forms.
            * `omit` - to omit the original tokens and include only the encoded tokens in the output of the token filter. Use this value if you want to only query on the encoded forms of the original tokens.
    * [lowercase](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-lowercase-tf-ref) - Normalizes token text to lowercase.
    * [length](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-length-tf-ref) - Removes tokens that are too short or too long:
        * `type` - Must be `length`
        * `min` - The minimum length of a token. Must be less than or equal to `max`.
        * `max` - The maximum length of a token. Must be greater than or equal to `min`.
    * [icuFolding](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-icufolding-tf-ref) - Applies character folding from [Unicode Technical Report #30](http://www.unicode.org/reports/tr30/tr30-4.html).
    * [icuNormalizer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-icunormalizer-tf-ref) - Normalizes tokens using a standard [Unicode Normalization Mode](https://unicode.org/reports/tr15/):
        * `type` - Must be 'icuNormalizer'.
        * `normalizationForm` - Normalization form to apply. Accepted values are:
            * `nfd` (Canonical Decomposition)
            * `nfc` (Canonical Decomposition, followed by Canonical Composition)
            * `nfkd` (Compatibility Decomposition)
            * `nfkc` (Compatibility Decomposition, followed by Canonical Composition)

      For more information about the supported normalization forms, see [Section 1.2: Normalization Forms, UTR#15](https://unicode.org/reports/tr15/#Norm_Forms).
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tf-ref) - Tokenizes input into n-grams of configured sizes.
        * `type` - Must be `nGram`
        * `minGram` - (Required) The minimum length of generated n-grams. Must be less than or equal to `maxGram`.
        * `maxGram` - (Required) The maximum length of generated n-grams. Must be greater than or equal to `minGram`.
        * `termNotInBounds` - Accepted values are:
            * `include`
            * `omit`

      If `include` is specified, tokens shorter than `minGram` or longer than `maxGram` are indexed as-is. If `omit` is specified, those tokens are not indexed.
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tf-ref) - Tokenizes input into edge n-grams of configured sizes:
        * `type` - Must be `edgeGram`
        * `minGram` - (Required) The minimum length of generated n-grams. Must be less than or equal to `max_gram`.
        * `maxGram` - (Required) The maximum length of generated n-grams. Must be greater than or equal to `min_gram`.
        * `termsNotInBounds` - Accepted values are:
            * `include`
            * `omit`

      If `include` is specified, tokens shorter than `minGram` or longer than `maxGram` are indexed as-is. If `omit` is specified, those tokens are not indexed.
    * [shingle](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-shingle-tf-ref) - Constructs shingles (token n-grams) from a series of tokens.
        * `type` - Must be `shingle`
        * `minShingleSize` - (Required) Minimum number of tokens per shingle. Must be less than or equal to `maxShingleSize`.
        * `maxShingleSize` - (Required) Maximum number of tokens per shingle. Must be greater than or equal to `minShingleSize`.
    * [regex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-regex-tf-ref) - Applies a regular expression to each token, replacing matches with a specified string.
        * `type` - Must be `regex`
        * `pattern` - (Required) Regular expression pattern to apply to each token.
        * `replacement` - (Required) Replacement string to substitute wherever a matching pattern occurs.
        * `matches` - (Required) Acceptable values are:
            * `all`
            * `first`

      If `matches` is set to `all, replace all matching patterns. Otherwise, replace only the first matching pattern.
    * [snowballStemming](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-snowballstemming-tf-ref) - Stems tokens using a [Snowball-generated stemmer](https://snowballstem.org/).
        * `type` - Must be `snowballstemming`
        * `stemmerName` - (Required) The following values are valid:
            * `arabic`
            * `armenian`
            * `basque`
            * `catalan`
            * `danish`
            * `dutch`
            * `english`
            * `finnish`
            * `french`
            * `german`
            * `german2` (Alternative German language stemmer. Handles the umlaut by expanding ü to ue in most contexts.)
            * `hungarian`
            * `irish`
            * `italian`
            * `kp` (Kraaij-Pohlmann stemmer, an alternative stemmer for Dutch.)
            * `lithuanian`
            * `lovins` (The first-ever published "Lovins JB" stemming algorithm.)
            * `norwegian`
            * `porter` (The original Porter English stemming algorithm.)
            * `portuguese`
            * `romanian`
            * `russian`
            * `spanish`
            * `swedish`
            * `turkish`
    * [stopword](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-stopword-tf-ref) - Removes tokens that correspond to the specified stop words. This token filter doesn't analyze the specified stop word:
        * `type` - Must be `stopword`
        * `token` - (Required) The list of stop words that correspond to the tokens to remove. Value must be one or more stop words.
        * `ignoreCase` - The flag that indicates whether or not to ignore case of stop words when filtering the tokens to remove. The value can be one of the following:
            * `true` - to ignore case and remove all tokens that match the specified stop words
            * `false` - to be case-sensitive and remove only tokens that exactly match the specified case

      If omitted, defaults to `true`.
    * [trim](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-trim-tf-ref) - Trims leading and trailing whitespace from tokens.

> **NOTE:** Do not use [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) token filter in operators where `fuzzy` is enabled. Atlas Search supports the `fuzzy` option for the following operators:
>* [autocomplete](https://docs.atlas.mongodb.com/reference/atlas-search/autocomplete/#std-label-autocomplete-ref)
>* [term (Deprecated)](https://docs.atlas.mongodb.com/reference/atlas-search/term/#std-label-term-ref)
>* [text](https://docs.atlas.mongodb.com/reference/atlas-search/text/#std-label-text-ref)



### Synonyms (search  index)
Synonyms mapping definition to use in the index.
* `name` - (Required) Name of the [synonym mapping definition](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-ref). Name must be unique in this index definition and it can't be an empty string.
* `source_collection` - (Required) Name of the source MongoDB collection for the synonyms. Documents in this collection must be in the format described in the [Synonyms Source Collection Documents](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-coll-spec).
* `analyzer` - (Required) Name of the [analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use with this synonym mapping. Atlas Search doesn't support these [custom analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) tokenizers and token filters in [analyzers used in synonym mappings](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#options):
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tokenizer-ref) Tokenizer
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tokenizer-ref) Tokenizers
    * [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) token filter
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tf-ref) token filter
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tf-ref) token filter
    * [shingle](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-shingle-tf-ref) token filter

```terraform
  synonyms {
   analyzer = "lucene.simple"
   name = "synonym_test"
   source_collection = "collection_test"
  }
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
