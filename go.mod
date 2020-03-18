module github.com/terraform-providers/terraform-provider-mongodbatlas

go 1.13

require (
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/beevik/etree v1.1.0 // indirect
	github.com/go-test/deep v1.0.1
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/hashicorp/terraform v0.12.1
	github.com/jen20/awspolicyequivalence v1.1.0 // indirect
	github.com/mongodb/go-client-mongodb-atlas v0.1.4-0.20200206183950-6d73c8e570f5
	github.com/mwielbut/pointy v1.1.0
	github.com/satori/uuid v1.2.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.0
	github.com/terraform-providers/terraform-provider-aws v1.9.0
	github.com/terraform-providers/terraform-provider-template v1.0.0 // indirect
	github.com/terraform-providers/terraform-provider-tls v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/net v0.0.0-20191009170851-d66e71096ffb // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	google.golang.org/grpc v1.23.0 // indirect
)

replace github.com/mongodb/go-client-mongodb-atlas => ../go-client-mongodb-atlas
