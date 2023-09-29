module github.com/zambien/terraform-provider-apigee

go 1.13

require (
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/hashicorp/terraform v0.12.17
	github.com/sethgrid/pester v0.0.0-20190127155807-68a33a018ad0 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/tibers/go-apigee-edge v0.0.0-20191119135131-525ca3781716
	github.com/zambien/go-apigee-edge v0.0.0-20191101145538-e45257f96262
	golang.org/x/crypto v0.6.0 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
