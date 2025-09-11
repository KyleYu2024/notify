module demo-plugin

go 1.24.2

require github.com/jianxcao/notify/backend v0.0.0-20250911025102-6c9cc29ce12a

replace github.com/jianxcao/notify/backend => ../../

require (
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
