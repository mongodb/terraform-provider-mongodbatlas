package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/examples-generation/prompts"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	// "log"
	// "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	// "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

const atlasAdminAPISpecURL = "https://raw.githubusercontent.com/mongodb/openapi/refs/heads/main/openapi/v2/openapi-2025-03-12.yaml"
const specFilePath = "open-api-spec.yml"

var resourceToGetPath = map[string]string{
	"alert_configuration": "/api/atlas/v2/groups/{groupId}/alertConfigs/{alertConfigId}",
	"search_index":        "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/search/indexes/{indexId}",
}

func main() {
	const resourceName = "alert_configuration"

	const DefaultAPIVersion = "2024-12-01-preview"
	const Model = "gpt-4.1"

	apiSpecSchema := getAPISpecSchema(resourceName)
	userPrompt := prompts.GetUserPrompt(prompts.UserPromptTemplateInputs{
		ResourceName:                  resourceName,
		ResourceImplementationSchema:  schemaContent,
		ResourceAPISpecResponseSchema: apiSpecSchema,
	})
	log.Printf("User Prompt: %s\n", userPrompt)

	client, err := CreateOpenAIClientWithKey(DefaultAPIVersion)
	if err != nil {
		panic(fmt.Errorf("failed to create OpenAI client: %w", err))
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompts.GenerateHCLSystemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model:       Model,
		Temperature: param.NewOpt(0.0),
	})

	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}

func getAPISpecSchema(resourceName string) string {
	getPath, ok := resourceToGetPath[resourceName]
	if !ok {
		return ""
	}

	if err := openapi.DownloadOpenAPISpec(atlasAdminAPISpecURL, specFilePath); err != nil {
		log.Fatalf("an error occurred when downloading Atlas Admin API spec: %v", err)
	}

	apiSpec, err := openapi.ParseAtlasAdminAPI(specFilePath)
	if err != nil {
		log.Fatalf("unable to parse Atlas Admin API: %v", err)
	}
	op, _ := apiSpec.Model.Paths.PathItems.Get(getPath)

	okResponse, _ := op.Get.Responses.Codes.Get(codespec.OASResponseCodeOK)
	schema, _ := codespec.GetSchemaFromMediaType(okResponse.Content)
	baseSchema := schema.Schema
	yamlBytes, _ := baseSchema.RenderInline()
	yamlString := string(yamlBytes)
	return yamlString
}

const schemaContent = `
func (r *alertConfigurationRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alert_configuration_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"matcher": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"field_name": schema.StringAttribute{
							Required: true,
						},
						"operator": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"metric_threshold_config": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"metric_name": schema.StringAttribute{
							Required: true,
						},
						"operator": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("GREATER_THAN", "LESS_THAN"),
							},
						},
						"threshold": schema.Float64Attribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Float64{
								float64planmodifier.UseStateForUnknown(),
							},
						},
						"units": schema.StringAttribute{
							Optional: true,
						},
						"mode": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
	}
}
`

const apiSpecSchemaContent = `
oneOf:
    - description: Other alerts which don't have extra details beside of basic one.
      properties:
        created:
            description: Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
            externalDocs:
                description: ISO 8601
                url: https://en.wikipedia.org/wiki/ISO_8601
            format: date-time
            readOnly: true
            type: string
        enabled:
            default: false
            description: Flag that indicates whether someone enabled this alert configuration for the specified project.
            type: boolean
`
