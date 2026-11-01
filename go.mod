module github.com/plusev-terminal/go-plugin-common/v2

go 1.24.2

require (
	github.com/extism/go-pdk v1.1.3
	github.com/go-playground/validator/v10 v10.28.0
	github.com/go-viper/mapstructure/v2 v2.4.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

// v2.0.0 added trading.Market.Extra, which breaks GORM AutoMigrate of models
// embedding trading.Market; v2.0.1 removed it again.
retract v2.0.0
