
# Development Best Practices

This document is the single source of truth for Terraform provider development best practices in this repository. Agent skills and agent configuration files (e.g. `AGENTS.md`, `CLAUDE.md`) should point to this document instead of duplicating its content.

## Table of Contents
- [Framework Choice](#framework-choice)
- [Schema Design](#schema-design)
  - [Optional vs Optional+Computed](#optional-vs-optionalcomputed)
  - [Keeping Computed-Only Fields on Resources](#keeping-computed-only-fields-on-resources)
  - [`UseStateForUnknown` Plan Modifier](#usestateforunknown-plan-modifier)
  - [`CreateOnly` Plan Modifier for Non-Updatable Attributes](#createonly-plan-modifier-for-non-updatable-attributes)
  - [Collections: Sets vs Lists](#collections-sets-vs-lists)
  - [Validation](#validation)
  - [SDK Getters](#sdk-getters)
- [Resource and Data Source Design](#resource-and-data-source-design)
  - [Pagination in Plural Data Sources](#pagination-in-plural-data-sources)
  - [Auto-Generated Data Source Schemas](#auto-generated-data-source-schemas)
  - [Configurable Timeouts](#configurable-timeouts)
- [State and Lifecycle](#state-and-lifecycle)
  - [Usage of `id`](#usage-of-id)
  - [Resources Deleted Outside Terraform](#resources-deleted-outside-terraform)
  - [Computed Attribute Defaults: `null` vs `""`](#computed-attribute-defaults-null-vs-)
- [Avoiding Breaking Changes](#avoiding-breaking-changes)
- [Diagnostics](#diagnostics)
- [Scaffolding Initial Code and File Structure](#scaffolding-initial-code-and-file-structure)
- [Auto-Generating Resources \& Data Sources (Internal tool)](#auto-generating-resources--data-sources-internal-tool)
  - [Customizing Generated Resources \& Data Sources](#customizing-generated-resources--data-sources)

## Framework Choice

Always use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) (TPF) for new resources and data sources. Never use the legacy SDKv2 for new code.

## Schema Design

### Optional vs Optional+Computed

- Use Optional+Computed for API properties that return a computed or default value when no value is provided in the request. If the API does not return a computed/default value, do not define the attribute as computed: this keeps the config as the source of truth (if the attribute is not present, the value is not set) and gives users awareness, through a non-empty plan, if the value is managed or modified externally.
- Be aware of the tradeoff: with Optional+Computed, when a user removes an attribute from their config Terraform does not generate a plan diff because it treats the server-side value as authoritative, so removals go undetected and the previous value silently persists. Only use Optional+Computed when the API genuinely returns a server-computed default for a field the user did not set.
- [IPA-111](https://mongodb.github.io/ipa/111) states that API properties should define a clear owner (client or server), which avoids Optional+Computed attributes altogether: attributes are either Required/Optional (client-owned) or Computed (server-owned).
- If a nested attribute is Optional+Computed, all its optional child attributes must also be defined as Optional+Computed to avoid infinite non-empty plans.

### Keeping Computed-Only Fields on Resources

- Default: keep all computed-only (server-owned) fields on resources. The autogen tooling supports this by default, and it matches other major provider conventions (AWS, GCP, Azure).
- Exclusion of attributes should be considered on a case-by-case basis and requires prior team alignment:
  - `advanced_cluster` is an explicit exception, driven by its specific context.
  - Effective fields are another exception: per team agreement they are exposed only through data sources to reduce plan verbosity.
- Where possible, push the upstream team to mark server-owned outputs in the API spec via the `x-xgen-server-computed-immutable` extension ([IPA-131](https://mongodb.github.io/ipa/131/#x-xgen-server-computed-immutable)) instead of relying on manual `immutable_computed` overrides in [`tools/codegen/config.yml`](../tools/codegen/config.yml).
- If the concern is plan verbosity, it can be reduced with `UseStateForUnknown` (see below), but only for computed values that do not change after creation.

### `UseStateForUnknown` Plan Modifier

As mentioned in the [Terraform documentation](https://developer.hashicorp.com/terraform/plugin/framework/resources/plan-modification#usestateforunknown), this plan modifier reduces the amount of `(known after apply)` messages for computed attributes that will not change their value, showing the value stored in the state instead.

- Use it only for genuinely immutable server-computed values, e.g. `created_at`, where showing `(known after apply)` would be misleading once the value is present in the state after the initial creation.
- Be cautious with `UseStateForUnknown` on Optional+Computed attributes: it can cause "inconsistent result after apply" errors when the API returns a value that differs from the prior state. If you are not certain the value never changes server-side when not defined by the user, omit it; showing `(known after apply)` during plan is safe.

### `CreateOnly` Plan Modifier for Non-Updatable Attributes

Use `customplanmodifier.CreateOnly()` ([`internal/common/customplanmodifier/create_only.go`](../internal/common/customplanmodifier/create_only.go)) on attributes that can only be set at creation time and cannot be updated. This generates a clear error during plan if a user tries to change the value, instead of silently failing or producing confusing API errors.

### Collections: Sets vs Lists

When a collection is needed, do not default to lists (indexable, ordered, can have repeated elements): sometimes sets (non-indexable, no order expectation, no repeated elements) make more sense. For example, if roles are unordered and cannot be repeated, a set is better than a list. Users can always write `tolist(roles)[0]` if needed, and this way they are explicit about their intentions and know they cannot rely on order.

### Validation

In general, prefer to skip validation in the Terraform provider and keep validation on the server side.

### SDK Getters

In general, prefer SDK getters over direct field access, e.g. `project.GetTags()` (see [PR #2135 discussion](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2135/files#r1560682050)).

## Resource and Data Source Design

- When creating a new resource or data source, consider the design principle of [representing a single API object](https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resources-should-represent-a-single-api-object): a new resource should correspond to a single set of create, read, delete, and optionally update methods.
- For the attribute schema of a resource or data source, consider the design principle of [aligning closely with the underlying API](https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api).
- As a consequence of the two previous statements, when supporting a new isolated GET operation, favor defining a new data source over expanding an existing resource or data source. This ensures existing resources and data sources do not gain attributes unrelated to the main API object. Example: [`mongodbatlas_control_plane_ip_addresses`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/data-sources/control_plane_ip_addresses).

### Pagination in Plural Data Sources

Pagination should not be supported in plural data sources (`items_per_page`, `page_num`, `include_count` and `total_count` may be ignored). Instead, the data source Read handler should iterate through all available data. Exceptions may apply if the corresponding Atlas resource could result in a huge number of records; in such cases pagination may be supported to avoid numerous calls to the API.

### Auto-Generated Data Source Schemas

Favor auto-generated singular and plural data source schemas from the resource schema using `DataSourceSchemaFromResource` and `PluralDataSourceSchemaFromResource` ([`internal/common/conversion/schema_generation.go`](../internal/common/conversion/schema_generation.go)).

### Configurable Timeouts

For resources where operations take some time before they transition to a desirable state, expose configurable [timeout arguments](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts) (example: [`internal/serviceapi/searchdeploymentapi/resource_schema.go`](../internal/serviceapi/searchdeploymentapi/resource_schema.go)). These must be included if there can be variation in the amount of time an operation takes, such as for the cluster resource where a creation can take 30 seconds or 1 hour depending on the provided `instance_size`. If a resource is consistent in the amount of time needed for an operation, configurable timeouts do not need to be defined given a proper default timeout is in place.

## State and Lifecycle

### Usage of `id`

The `id` concept appears in different usage contexts (see the [resource lifecycle](https://github.com/hashicorp/terraform/blob/main/docs/resource-instance-change-lifecycle.md) for background):

1. **Import ID**: usually a string with `-` separating the id parts, e.g. `{project_id}-{instance_name}` for `mongodbatlas_stream_instance` ([`internal/service/streaminstance/resource_stream_instance.go`](../internal/service/streaminstance/resource_stream_instance.go)).
2. **Internal resource identification in SDKv2**: [`SetId()`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ResourceData.SetId) sets the current resource id (we usually use `EncodeStateID` to base64-encode a map of attributes) and [`Id()`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ResourceData.Id) reads it back (`DecodeStateID` parses the map back, see [`internal/common/conversion/encode_state.go`](../internal/common/conversion/encode_state.go)). Recommendation: do not expose this field to users.
3. **Attribute `id`**: in TPF it is OK to use `id` as an attribute, e.g. [`stream_connection`](../internal/service/streamconnection/resource_stream_connection.go) and [`stream_instance`](../internal/service/streaminstance/resource_stream_instance.go). In SDKv2, do not use it as an attribute: it is confusing to have both an attribute named `id` and the internal id field, e.g. `d.Id()` reads the internal id while `d.Get("id").(string)` reads the attribute set by the user.

### Resources Deleted Outside Terraform

A resource needs to be removed from the state when it has been deleted outside of Terraform (see [PR #2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268)):

- SDKv2: `d.SetId("")`
- TPF: `resp.State.RemoveResource(ctx)`

### Computed Attribute Defaults: `null` vs `""`

When an attribute has been removed from the API but we want to avoid breaking clients, we must understand the "empty value":

1. The attribute has been set to `""` (empty string), usually the case if we use the auto-generated SDK with `GetXXX()` getters. In this case, we must explicitly set it to `""` in the provider to avoid changes to the state:
   - SDKv2: `d.Set("attr", "")` ([example for `err_msg`](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2255/commits/5d2ce594aa91e730063ff18b4f7fac9ce0065ca1))
   - TPF: set `Attr: ""` when creating the model from the API payload
2. The attribute does not exist in the state, aka `null`. No need to set anything in the state, both SDKv2 and TPF will use `null` as the default.

## Avoiding Breaking Changes

Any change to a resource or data source must be **backward-compatible by default**. Treat the following as breaking unless the intent is explicitly documented and acknowledged:

- **Schema changes**: renaming, removing, or changing the type of an existing attribute; converting between Optional, Required, and Computed; tightening or removing valid values from a validation.
- **State representation**: altering how a value is stored in or read from state (e.g. changing casing, encoding, flattening/nesting structure). Even if the API value is the same, a different state representation triggers unexpected diffs or "inconsistent result after apply" errors.
- **Default value changes**: adding, removing, or modifying a `Default` or plan modifier that changes what gets written to state when the user omits a field.
- **Import behavior**: changing the format of the import ID or what is read during import.

When a breaking change is genuinely required (e.g. aligning with an API deprecation), call it out explicitly in the PR description.

## Diagnostics

Per the [HashiCorp diagnostics guidance](https://developer.hashicorp.com/terraform/plugin/framework/diagnostics):

- `Summary` should be a real sentence, static, and unique within the resource. Unique and static makes it easy to debug by searching for the text and going directly to the code line.
- `Detail` can be more sentences with parameters. Tell the practitioner exactly what they need to fix and how.
- Use `Warning` to inform practitioners about suboptimal situations they should resolve to ensure stable functioning (e.g. deprecations), or to inform them about possible unexpected behaviors.

## Scaffolding Initial Code and File Structure

**Note**: This command is relevant when developing a new resource or data source manually. For full autogeneration, reference the section below.

This command can be used the following way:
```bash
make scaffold resource_name=streamInstance type=resource
```
- **resource_name**: The name of the resource, which must be defined in camel case.
- **type**: Describes the type of resource being created. There are 3 different types: `resource`, `data-source`, `plural-data-source`.

This will generate resource/data source files and accompanying test files needed for starting the development, and will contain multiple comments with `TODO:` statements which give guidance for the development.

As a follow up step, use [Auto-Generating Resources](#auto-generating-resources--data-sources-internal-tool) to autogenerate the schema via the Open API specification. This will require making adjustments to the generated `./internal/service/<resource_name>/tfplugingen/generator_config.yml` file.

## Auto-Generating Resources & Data Sources (Internal tool)

The generation command makes use of a configuration file defined under [`./tools/codegen/config.yml`](../tools/codegen/config.yml). The structure of this configuration file can be found under  [`./tools/codegen/config/config_model.go`](../tools/codegen/config/config_model.go).

The generation command takes a single optional argument `resource_name`. If not provided, all resources defined in the configuration are generated.

```bash
make autogen-pipeline resource_name=search_deployment_api
```

If you wish to generate resource/data source models without fetching latest changes from the API Spec, use the following command:
```bash
make autogen-generate-resources
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

### Acceptance tests

While no form of acceptance test generation is available, it is expected for autogenerated resource to have acceptance tests. These can be defined directly in:
`./internal/serviceapi/<resource-package>/resource_test.go` 

## Customizing Generated Resources & Data Sources

There are two primary ways to customize the behavior and schema of autogenerated resources and data sources: **configuration-driven** options in `tools/codegen/config.yml` and **custom code hooks** implemented on the generated types.

### Configuration-Driven Customizations (`tools/codegen/config.yml`)

These options are defined per resource in `tools/codegen/config.yml` (see the configuration model in `tools/codegen/config/config_model.go` and `tools/codegen/codespec/config.go`) and allow you to adjust how schemas and docs are generated:

- **`ignores`**: Hide attributes that are not useful to end users.

- **`aliases`**: Aligns attribute names with existing provider conventions so configurations feel consistent across manual and autogenerated resources. For example:
  - `group_id` → `project_id` (defined in config file with camelCase: `groupId: projectId`)

- **Computability and requiredness overrides**: Use `schema.overrides.<attribute>.computability` to override whether an attribute is **required**, **optional**, or **computed** to better match API behavior. This includes support for attributes that are both **optional + computed**, which is useful for fields that may be omitted by users but are defaulted or populated by the API (for example, defaulted fields on clusters).

- **Description overrides**: Use `schema.overrides.<attribute>.description` sparingly to correct or clarify attribute descriptions when the OpenAPI or Atlas documentation is misleading, incomplete, or too low-level. Keep wording consistent with other resources and focused on what the Terraform user needs to know, without changing the underlying API semantics.

- **`sensitive`**: Marks attributes as sensitive, which is especially helpful when the API specification does not annotate secret properties with `format: password`. Use this to ensure secrets are not displayed in plan/apply outputs.

- **Type overrides for list/set types**: You can override the default type (e.g. `list` vs `set`) to better match the behavior of the underlying API. This is only supported for list/set types and should be used when ordering or uniqueness semantics matter to users.

When adding or changing these options, keep changes minimal, documented, and consistent with existing patterns in the same resource family.

### Custom Code Hooks (`internal/common/autogen/custom_hooks.go`)

For cases where configuration alone is not sufficient, autogenerated resources and data sources can implement hook interfaces defined in `internal/common/autogen/custom_hooks.go`. These hooks allow you to intercept and customize behavior around resource/data source logic while still relying on the default generator for the bulk of the implementation.

Available hooks include:

- **`PreReadAPICallHook` / `PostReadAPICallHook`** (relevant for both resource and data sources)
- **`PreCreateAPICallHook` / `PostCreateAPICallHook`**
- **`PreUpdateAPICallHook` / `PostUpdateAPICallHook`**
- **`PreDeleteAPICallHook` / `PostDeleteAPICallHook`**

To use them, implement the relevant interface(s) on the autogenerated resource or data source struct in the corresponding `internal/serviceapi/<resource-package>/` package. Resource/data source implementation will detect and invoke these functions if available.

**Note:** In order to preserve automatic updates in autogenerated code, it is advised custom code hooks do not leverage typed Go SDK models. Custom code should be future proof to supporting addition of new attributes in the schema without any changes in hook implementations.

An example implementation can be found in:

- `internal/serviceapi/orgserviceaccountsecretapi/resource_custom_hooks.go`

In this example, a custom `PostReadAPICall` implementation filters the raw API response to return only the specific secret matching the Terraform resource’s ID, mimicking the read response of a single secret which is not available in the API.

