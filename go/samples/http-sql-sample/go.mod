module http-sql-sample

go 1.19

require (
	github.com/go-sql-driver/mysql v1.6.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	google.com/sqlcommenter/http-net v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20220927170352-d9d178bc13c6 // indirect
	google.com/sqlcommenter/core v0.0.0-00010101000000-000000000000
	google.com/sqlcommenter/go-sql v0.0.0-00010101000000-000000000000
)

// replace your local path if choosing local install

//replace google.com/sqlcommenter/core => /Users/thiyagunataraj/Documents/Sqlcommenter-new/sqlcommenter/go/core

//replace google.com/sqlcommenter/go-sql => /Users/thiyagunataraj/Documents/Sqlcommenter-new/sqlcommenter/go/go-sql

//replace google.com/sqlcommenter/http-net => /Users/thiyagunataraj/Documents/Sqlcommenter-new/sqlcommenter/go/http-net
