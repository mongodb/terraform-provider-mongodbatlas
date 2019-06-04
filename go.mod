module github.com/terraform-providers/terraform-provider-mongodbatlas

go 1.12

require (
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/hashicorp/terraform v0.12.1
	github.com/mongodb-partners/go-client-mongodb-atlas v0.0.0
)

replace github.com/mongodb-partners/go-client-mongodb-atlas v0.0.0 => ../go-client-mongodb-atlas/
