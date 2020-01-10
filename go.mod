module github.com/terraform-providers/terraform-provider-mongodbatlas

go 1.12

require (
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce // indirect
	github.com/hashicorp/terraform v0.12.1
	github.com/mongodb/go-client-mongodb-atlas v0.1.2
	github.com/mwielbut/pointy v1.1.0
	github.com/spf13/cast v1.3.0
)

replace github.com/mongodb/go-client-mongodb-atlas v0.1.2 => ../go-client-mongodb-atlas
