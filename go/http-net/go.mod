module github.com/google/sqlcommenter/go/http-net

replace google.com/sqlcommenter/core => ../core/.

go 1.19

require (
	github.com/google/sqlcommenter/go/core v0.0.1-beta
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
)
