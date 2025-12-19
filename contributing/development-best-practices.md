
# Development Best Practices

## Table of Contents
- [Development Best Practices](#development-best-practices)
  - [Table of Contents](#table-of-contents)
  - [Scaffolding Initial Code and File Structure](#scaffolding-initial-code-and-file-structure)
  - [Auto-Generating Resources \& Data Sources](#auto-generating-resources--data-sources)
    - [(Recommended) Using internal tool](#recommended-using-internal-tool)

A set of commands have been defined with the intention of speeding up development process, while also preserving common conventions throughout our codebase.

## Scaffolding Initial Code and File Structure

This command can be used the following way:
```bash
make scaffold resource_name=streamInstance type=resource
```
- **resource_name**: The name of the resource, which must be defined in camel case.
- **type**: Describes the type of resource being created. There are 3 different types: `resource`, `data-source`, `plural-data-source`.

This will generate resource/data source files and accompanying test files needed for starting the development, and will contain multiple comments with `TODO:` statements which give guidance for the development.

As a follow up step, use [Auto-Generating Resources](#auto-generating-resources) to autogenerate the schema via the Open API specification. This will require making adjustments to the generated `./internal/service/<resource_name>/tfplugingen/generator_config.yml` file.

## Auto-Generating Resources & Data Sources

### (Recommended) Using internal tool

The generation command makes use of a configuration file defined under [`./tools/codegen/config.yml`](../tools/codegen/config.yml). The structure of this configuration file can be found under  [`./tools/codegen/config/config_model.go`](../tools/codegen/config/config_model.go).

The generation command takes a single optional argument `resource_name`. If not provided, all resources defined in the configuration are generated.

```bash
make resource-generation-pipeline resource_name=search_deployment_api
```


As a result, content of schemas and models will be written into the corresponding resource packages:
`./internal/serviceapi/<resource-package>/resource_schema.go`

And operations will be written into:
`./internal/serviceapi/<resource-package>/resource.go`

Data sources are automatically generated as part of the same generation process when a `datasources` block is configured in `tools/codegen/config.yml`. The tool generates both singular and plural data sources:

**Singular Data Source** (generated when `datasources.read` is configured):
- `./internal/serviceapi/<resource-package>/data_source_schema.go`
- `./internal/serviceapi/<resource-package>/data_source.go`

**Plural Data Source** (generated when `datasources.list` is configured):
- `./internal/serviceapi/<resource-package>/plural_data_source_schema.go`
- `./internal/serviceapi/<resource-package>/plural_data_source.go`

