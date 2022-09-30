module google.com/sqlcommenter/gosql

go 1.19

require (
	go.opentelemetry.io/otel/sdk v1.10.0
	google.com/sqlcommenter/core v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20220927170352-d9d178bc13c6 // indirect
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	google.com/sqlcommenter/http-net v0.0.0-00010101000000-000000000000
)

replace google.com/sqlcommenter/core => ../core/.

replace google.com/sqlcommenter/http-net => ../http-net/.
