package sources

type DataflowSourceConfig struct {
	Project  string `validate:"empty=false" yaml:"project"`
	Location string `validate:"empty=false" yaml:"location"`
}
