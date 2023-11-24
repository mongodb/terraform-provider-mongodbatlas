package config

const (
	ProviderConfigError                      = "error in configuring the provider."
	AWS                                      = "AWS"
	AZURE                                    = "AZURE"
	MongodbGovCloudURL                       = "https://cloud.mongodbgov.com"
	ToolName                                 = "terraform-provider-mongodbatlas"
	MissingAuthAttrError                     = "either Atlas Programmatic API Keys or AWS Secrets Manager attributes must be set"
	DeprecationParamByDate                   = "this parameter is deprecated and will be removed by %s"
	DeprecationParamByDateWithReplacement    = "this parameter is deprecated and will be removed by %s, please transition to %s"
	DeprecationParamByVersion                = "this parameter is deprecated and will be removed in version %s"
	DeprecationResourceByDateWithReplacement = "this resource is deprecated and will be removed in %s, please transition to %s"
	ErrorProjectSetting                      = "error setting `%s` for project (%s): %s"
	ErrorGetRead                             = "error reading cloud provider access %s"
)
