package signmw

const configSectionName = "auth"

type signatureConfig struct {
	RequestSigningKey string `yaml:"request_signing_key"`
}
