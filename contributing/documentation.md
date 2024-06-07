# Documentation

- In our documentation, when a resource field allows a maximum of only one item, we do not format that field as an array. Instead, we create a subsection specifically for this field. Within this new subsection, we enumerate all the attributes of the field. Let's illustrate this with an example: [cloud_backup_schedule.html.markdown](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/r/cloud_backup_schedule.html.markdown?plain=1#L207)
- You can check how the documentation is rendered on the Terraform Registry via [doc-preview](https://registry.terraform.io/tools/doc-preview).

## Creating Resource and Data source Documentation
We autogenerate the documentation of our provider resources and data sources via [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

### How to generate the documentation for a resource
- Make sure that the resource and data source schemas have defined the fields `MarkdownDescription` and `Description`.
  - We recommend to use [Scaffolding Schema and Model Definitions](#scaffolding-schema-and-model-definitions) to autogenerate the schema via the Open API specification.
- Add the resource/data source templates to the [templates](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/templates) folder. See [README.md](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/templates/README.md) for more info.
- Run the Makefile command `generate-doc`
```bash
export resource_name=search_deployment && make generate-doc
```
