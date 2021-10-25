---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: search index"
sidebar_current: "docs-mongodbatlas-resource-search-index"
description: |-
Provides a Search Index resource.
---

# mongodbatlas_search_index

`mongodbatlas_search_index` provides a Search Index resource. This allows indexes to be created.

## Example Usage

### Basic 
```hcl
resource "mongodbatlas_search_index" "test" {
  name   = "project-name"
  project_id = "<PROJECT_ID>"
  cluster_name = "<CLUSTER_NAME>"
  
  analyzer = "lucene.standard"
  collectionName = "collection_test"
  database = "database_test"
  mappings_dynamic = true
  
  searchAnalyzer = "lucene.standard"
}
```

### Advanced (with custom analyzers)
```hcl
resource "mongodbatlas_search_index" "test" {
  project_id = "%[1]s"
  cluster_name = "%[2]s"
  analyzer = "lucene.standard"
  collectionName = "collection_test"
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
  name = "name_test"
  searchAnalyzer = "lucene.standard"
  analyzers = <<-EOF
  [{
  "name": "index_analyzer_test_name",
  "char_filters": {
	"type": "mapping",
	"mappings": {"\\" : "/"}
    	},
  "tokenizer": {
  "type": "nGram",
  "min_gram": 2,
  "max_gram": 5
		},
  "token_filters": {
	"type": "length",
	"min": 20,
	"max": 33
    	}
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

* `name` - (Required) The name of the search index you want to create.
* `project_id` - (Required) The ID of the organization or project you want to create the search index within.
* `cluster_name` - (Required) The name of the cluster where you want to create the search index within.


* `analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when creating the index. Defaults to [lucene.standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/standard/#std-label-ref-standard-analyzer)

* `analyzers` - [Custom analyzers](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) to use in this index. This is an array of JSON objects.
```
analyzers = <<-EOF
  [{
  "name": "index_analyzer_test_name",
  "char_filters": {
	"type": "mapping",
	"mappings": {"\\" : "/"}
    	},
  "tokenizer": {
  "type": "nGram",
  "min_gram": 2,
  "max_gram": 5
	},
  "token_filters": {
	"type": "length",
	"min": 20,
	"max": 33
    	}
  }]
EOF
```

* `collection_name` - (Required) Name of the collection the index is on.

* `database` - (Required) Name of the database the collection is in.

* `mappings_dynamic` - Indicates whether the index uses dynamic or static mapping. For dynamic mapping, set the value to `true`. For static mapping, specify the fields to index using `mappings_fields`

* `mappings_fields` - attribute is required when `mappings_dynamic` is true. This field needs to be a JSON string in order to be decoded correctly.
  ```hcl
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
  ```

* `search_analyzer` - [Analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use when searching the index. Defaults to [lucene.standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/standard/#std-label-ref-standard-analyzer)
* `synonyms` - Synonyms mapping definition to use in this index.

### Analyzers
An [Atlas Search analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/) prepares a set of documents to be indexed by performing a series of operations to transform, filter, and group sequences of characters. You can define a custom analyzer to suit your specific indexing needs.

* `name` - (Required) 	
  Name of the custom analyzer. Names must be unique within an index, and may **not** start with any of the following strings:
    * `lucene`
    * `builtin`
    * `mongodb`
* `char_filters` - Array containing zero or more character filters. Always require a `type` field, and some take additional options as well
  ```hcl
  "char_filters":{
   "type": "<FILTER_TYPE>",
   "ADDITIONAL_OPTION": VALUE
  }
  ```
  Atlas search supports four `types` of character filters:
    * [htmlStrip](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-htmlStrip-ref) - Strips out HTML constructs
        * `type` - (Required) Must be `htmlStrip`
        * `ignored_tags`- a list of HTML tags to exclude from filtering
        ```hcl
          analyzers = <<-EOF [{
            "name": "analyzer_test",
            "char_filters":{
              "type": "htmlStrip",
              "ignored_tags": ["a"]
              }   
            }] 
       ```
    * [icuNormalize](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-icuNormalize-ref) - Normalizes text with the [ICU](http://site.icu-project.org/) Normalizer. Based on Lucene's [ICUNormalizer2CharFilter](https://lucene.apache.org/core/8_3_0/analyzers-icu/org/apache/lucene/analysis/icu/ICUNormalizer2CharFilter.html)
    * [mapping](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-mapping-ref) - Applies user-specified normalization mappings to characters. Based on Lucene's [MappingCharFilter](https://lucene.apache.org/core/8_0_0/analyzers-common/org/apache/lucene/analysis/charfilter/MappingCharFilter.html)
      An object containing a comma-separated list of mappings. A mapping indicates that one character or group of characters should be substituted for another, in the format `<original> : <replacement>`
      ### Example
        ```hcl
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


* `tokenizer` - (Required) Tokenizer to use. Determines how Atlas Search splits up text into discrete chunks of indexing. Always require a type field, and some take additional options as well.
    ```hcl
    "tokenizer":{
    "type": "<tokenizer-type>",
    "ADDITIONAL_OPTIONS": VALUE
    }
    ```
  Atlas Search supports the following tokenizer options:
    * [standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-standard-tokenizer-ref) - Tokenize based on word break rules from the [Unicode Text Segmentation algorithm](http://www.unicode.org/L2/L2019/19034-uax29-34-draft.pdf):
        * `type` - Must be `standard`
        * `max_token_length` - Maximum length for a single token. Tokens greater than this length are split at `maxTokenLength` into multiple tokens.
    * [keyword](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-keyword-tokenizer-ref) - Tokenize the entire input as a single token.
    * [whitespace](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-whitespace-tokenizer-ref) - Tokenize based on occurrences of whitespace between words.
        * `type` - Must be `whitespace`
        * `max_token_length` - Maximum length for a single token. Tokens greater than this length are split at `maxTokenLength` into multiple tokens.
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tokenizer-ref) - Tokenize into text chunks, or "n-grams", of given sizes.
        * `type` - Must be `nGram`
        * `min_gram` - (Required) Number of characters to include in the shortest token created.
        * `max_gram` - (Required) Number of characters to include in the longest token created.
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tokenizer-ref) - Tokenize input from the beginning, or "edge", of a text input into n-grams of given sizes.
        * `type` - Must be `edgeGram`
        * `min_gram` - (Required) Number of characters to include in the shortest token created.
        * `max_gram` - (Required) Number of characters to include in the longest token created.
    * [regexCaptureGroup](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-regexcapturegroup-tokenizer-ref) - Match a regular expression pattern to extract tokens.
        * `type` - Must be `regexCaptureGroup`
        * `pattern` - (Required) A regular expression to match against.
        * `group` - (Required) Index of the character group within the matching expression to extract into tokens. Use 0 to extract all character groups.
    * [regexSplit](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-regexSplit-tokenizer-ref) - Split tokens with a regular-expression based delimiter.
        * `type` - Must be `regexSplit`
        * `pattern` - (Required) A regular expression to match against.
    * [uaxUrlEmail](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-uaxUrlEmail-tokenizer-ref) - Tokenize URLs and email addresses. Although [uaxUrlEmail](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-uaxUrlEmail-tokenizer-ref) tokenizer tokenizes based on word break rules from the [Unicode Text Segmentation algorithm](http://www.unicode.org/L2/L2019/19034-uax29-34-draft.pdf), we recommend using uaxUrlEmail tokenizer only when the indexed field value includes URLs and email addresses. For fields that do not include URLs or email addresses, use the [standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-standard-tokenizer-ref) tokenizer to create tokens based on word break rules.
        * `type` - Must be `uaxUrlEmail`
        *  `max_token_length` - The maximum number of characters in one token.

* `tokenFilters` - Array containing zero or more token filters. Always require a type field, and some take additional options as well:
  ```hcl
  "token_filters":{
    "type": "<FILTER_TYPE>",
    "ADDITIONAL-OPTIONS": VALUE
  }
  ```
  Atlas Search supports the following token filters:
    * [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) - Creates tokens for words that sound the same based on [Daitch-Mokotoff Soundex](https://en.wikipedia.org/wiki/Daitch%E2%80%93Mokotoff_Soundex) phonetic algorithm. This filter can generate multiple encodings for each input, where each encoded token is a 6 digit number:
        *  `type` - Must be `daitchMokotoffSoundex`
        * `original_tokens` - Specifies whether to include or omit the original tokens in the output of the token filter. Value can be one of the following:
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
        * `normalization_form` - Normalization form to apply. Accepted values are:
            * `nfd` (Canonical Decomposition)
            * `nfc` (Canonical Decomposition, followed by Canonical Composition)
            * `nfkd` (Compatibility Decomposition)
            * `nfkc` (Compatibility Decomposition, followed by Canonical Composition)

      For more information about the supported normalization forms, see [Section 1.2: Normalization Forms, UTR#15](https://unicode.org/reports/tr15/#Norm_Forms).
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tf-ref) - Tokenizes input into n-grams of configured sizes.
        * `type` - Must be `nGram`
        * `min_gram` - (Required) The minimum length of generated n-grams. Must be less than or equal to `maxGram`.
        * `max_gram` - (Required) The maximum length of generated n-grams. Must be greater than or equal to `minGram`.
        * `terms_not_in_bounds` - Accepted values are:
            * `include`
            * `omit`

      If `include` is specified, tokens shorter than `min_gram` or longer than `max_gram` are indexed as-is. If `omit` is specified, those tokens are not indexed.
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tf-ref) - Tokenizes input into edge n-grams of configured sizes:
        * `type` - Must be `edgeGram`
        * `min_gram` - (Required) The minimum length of generated n-grams. Must be less than or equal to `max_gram`.
        * `max_gram` - (Required) The maximum length of generated n-grams. Must be greater than or equal to `min_gram`.
        * `terms_not_in_bounds` - Accepted values are:
            * `include`
            * `omit`

      If `include` is specified, tokens shorter than `min_gram` or longer than `max_gram` are indexed as-is. If `omit` is specified, those tokens are not indexed.
    * [shingle](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-shingle-tf-ref) - Constructs shingles (token n-grams) from a series of tokens.
        * `type` - Must be `shingle`
        * `min_shingle_size` - (Required) Minimum number of tokens per shingle. Must be less than or equal to `max_shingle_size`.
        * `max_shingle_size` - (Required) Maximum number of tokens per shingle. Must be greater than or equal to `min_shingle_size`.
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
        * `stemmer_name` - (Required) The following values are valid:
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
        * `ignore_case` - The flag that indicates whether or not to ignore case of stop words when filtering the tokens to remove. The value can be one of the following:
            * `true` - to ignore case and remove all tokens that match the specified stop words
            * `false` - to be case-sensitive and remove only tokens that exactly match the specified case

      If omitted, defaults to `true`.
    * [trim](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-trim-tf-ref) - Trims leading and trailing whitespace from tokens.

> **NOTE:** Do not use [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) token filter in operators where `fuzzy` is enabled. Atlas Search supports the `fuzzy` option for the following operators:
>* [autocomplete](https://docs.atlas.mongodb.com/reference/atlas-search/autocomplete/#std-label-autocomplete-ref)
>* [term (Deprecated)](https://docs.atlas.mongodb.com/reference/atlas-search/term/#std-label-term-ref)
>* [text](https://docs.atlas.mongodb.com/reference/atlas-search/text/#std-label-text-ref)



### Synonyms
Synonyms mapping definition to use in the index.
* `name` - (Required) Name of the [synonym mapping definition](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-ref). Name must be unique in this index definition and it can't be an empty string.
* `source_collection` - (Required) Name of the source MongoDB collection for the synonyms. Documents in this collection must be in the format described in the [Synonyms Source Collection Documents](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-coll-spec).
* `analyzer` - (Required) Name of the [analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/#std-label-analyzers-ref) to use with this synonym mapping. If you set `mappings.dynamic` to `true`, you must use the default analyzer, [lucene.standard](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/standard/#std-label-ref-standard-analyzer), here also. Atlas Search doesn't support these [custom analyzer](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-custom-analyzers) tokenizers and token filters in the index definition for [synonyms](https://docs.atlas.mongodb.com/reference/atlas-search/synonyms/#std-label-synonyms-ref):
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tokenizer-ref) Tokenizer
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tokenizer-ref) Tokenizers
    * [daitchMokotoffSoundex](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-daitchmokotoffsoundex-tf-ref) token filter
    * [nGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-ngram-tf-ref) token filter
    * [edgeGram](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-edgegram-tf-ref) token filter
    * [shingle](https://docs.atlas.mongodb.com/reference/atlas-search/analyzers/custom/#std-label-shingle-tf-ref) token filter

```hcl
  synonyms {
   analyzer = "lucene.simple"
   name = "synonym_test"
   source_collection = "collection_test"
  }
```



For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/atlas-search/) - [and MongoDB Atlas API - Search](https://docs.atlas.mongodb.com/reference/api/atlas-search/) Documentation for more information.
