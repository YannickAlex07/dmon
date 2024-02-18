package slack

type SlackService interface {
	Send() error
}
