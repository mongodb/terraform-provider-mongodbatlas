# Templates 
This folder contains the template files used by [TFplugindocs](https://github.com/hashicorp/terraform-plugin-docs) to autogenerate our provider documentation.


## How To Guide

### How do we use templates

The templates in [TFplugindocs](https://github.com/hashicorp/terraform-plugin-docs) are implemented with Go [`text/template`](https://golang.org/pkg/text/template/). After running `tfplugindocs generate`, the tfplugindocs engine performs the following operations:
 
 - Retrieves the templates in `docs/templates/resources/{resource_name}.html.markdown.tmpl` and `docs/templates/data-sources/{resource_name}.html.markdown.tmpl` for a specific `resource-name`
 - Retrieves the values in `MarkdownDescription` for all the fields in the `resource-name` schema
 - Retrieves the resource examples in `examples/{resource_name}/main.tf`
 - Generates the documentation.

 
 ### How to generate a new template
 Use [resources/search_deployment.html.markdown.tmpl](resources/search_deployment.html.markdown.tmpl) and [resources/search_deployment.html.markdown.tmpl](data-sources/search_deployment.html.markdown.tmpl) as an example to add templates for a resource and data source. 

#### Data fields
Here a list of the basic data fields and functions that you can use when defining a template.

See [the HashiCorp documentation](https://github.com/hashicorp/terraform-plugin-docs?tab=readme-ov-file#templates) for a full list of data fields and functions.

##### Provider

|                   Field |  Type  | Description                                                                               |
|------------------------:|:------:|-------------------------------------------------------------------------------------------|
|          `.Description` | string | Provider description                                                                      |
|           `.HasExample` |  bool  | Is there an example file?                                                                 |
|          `.ExampleFile` | string | Path to the file with the terraform configuration example                                 |
|         `.ProviderName` | string | Canonical provider name (ex. `terraform-provider-random`)                                 |
|    `.ProviderShortName` | string | Short version of the provider name (ex. `random`)                                         |
| `.RenderedProviderName` | string | Value provided via argument `--rendered-provider-name`, otherwise same as `.ProviderName` |
|       `.SchemaMarkdown` | string | a Markdown formatted Provider Schema definition                                           |

##### Resources / Data Source

|                   Field |  Type  | Description                                                                               |
|------------------------:|:------:|-------------------------------------------------------------------------------------------|
|                 `.Name` | string | Name of the resource/data-source (ex. `tls_certificate`)                                  |
|                 `.Type` | string | Either `Resource` or `Data Source`                                                        |
|          `.Description` | string | Resource / Data Source description                                                        |
|           `.HasExample` |  bool  | Is there an example file?                                                                 |
|          `.ExampleFile` | string | Path to the file with the terraform configuration example                                 |
|            `.HasImport` |  bool  | Is there an import file?                                                                  |
|           `.ImportFile` | string | Path to the file with the command for importing the resource                              |
|         `.ProviderName` | string | Canonical provider name (ex. `terraform-provider-random`)                                 |
|    `.ProviderShortName` | string | Short version of the provider name (ex. `random`)                                         |
| `.RenderedProviderName` | string | Value provided via argument `--rendered-provider-name`, otherwise same as `.ProviderName` |
|       `.SchemaMarkdown` | string | a Markdown formatted Resource / Data Source Schema definition                             |

#### Functions

| Function        | Description                                                                                       |
|-----------------|---------------------------------------------------------------------------------------------------|
| `codefile`      | Create a Markdown code block with the content of a file. Path is relative to the repository root. |
| `lower`         | Equivalent to [`strings.ToLower`](https://pkg.go.dev/strings#ToLower).                            |
| `plainmarkdown` | Render Markdown content as plaintext.                                                             |
| `prefixlines`   | Add a prefix to all (newline-separated) lines in a string.                                        |
| `printf`        | Equivalent to [`fmt.Printf`](https://pkg.go.dev/fmt#Printf).                                      |
| `split`         | Split string into sub-strings, by a given separator (ex. `split .Name "_"`).                      |
| `title`         | Equivalent to [`cases.Title`](https://pkg.go.dev/golang.org/x/text/cases#Title).                  |
| `tffile`        | A special case of the `codefile` function, designed for Terraform files (i.e. `.tf`).             |
| `trimspace`     | Equivalent to [`strings.TrimSpace`](https://pkg.go.dev/strings#TrimSpace).                        |
| `upper`         | Equivalent to [`strings.ToUpper`](https://pkg.go.dev/strings#ToUpper).                            |

