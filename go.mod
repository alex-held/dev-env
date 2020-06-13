module github.com/alex-held/dev-env

go 1.14

require (
	github.com/ebuchman/go-shell-pipes v0.0.0-20150412091402-83e132480862
	github.com/magiconair/properties v1.8.0
	github.com/rs/zerolog v1.19.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.10.0
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
)

replace github.com/alex-held/dev-env/utils => ./utils
