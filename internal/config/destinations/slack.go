package destinations

type SlackDestinationConfig struct {
	Token   string `validate:"empty=false" yaml:"token"`
	Channel string `validate:"empty=false" yaml:"channel"`
}
