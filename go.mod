module thingdust

go 1.18

require (
	github.com/eliona-smart-building-assistant/go-eliona v1.9.4
	github.com/eliona-smart-building-assistant/go-eliona-api-client/v2 v2.4.8
	github.com/eliona-smart-building-assistant/go-utils v1.0.17
	github.com/friendsofgo/errors v0.9.2
	github.com/gorilla/mux v1.8.0
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/sqlboiler/v4 v4.14.1
	github.com/volatiletech/strmangle v0.0.4
)

// Bugfix see: https://github.com/volatiletech/sqlboiler/blob/91c4f335dd886d95b03857aceaf17507c46f9ec5/README.md
// decimal library showing errors like: pq: encode: unknown type types.NullDecimal is a result of a too-new and broken version of the github.com/ericlargergren/decimal package, use the following version in your go.mod: github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5

require (
	github.com/ericlagergren/decimal v0.0.0-20221120152707-495c53812d05 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/randomize v0.0.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
)
