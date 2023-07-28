//go:generate oapi-codegen -old-config-style -generate gin -o openapi/openapi_gin.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -old-config-style -generate spec -o openapi/openapi_spec.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -old-config-style -generate types -o openapi/openapi_types.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -old-config-style -generate client -o openapi/openapi_client.gen.go -package openapi ../../spec/openapi.yaml

package api
