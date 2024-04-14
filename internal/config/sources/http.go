package sources

type HTTPSourceConfig struct {
	URL            string `validate:"empty=false & format=url" yaml:"url"`
	ExpectedStatus int    `yaml:"expected_status"`
}
