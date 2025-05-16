package domain

type (
	ChatType int

	MessageType   int
	MessageStatus int
)

const (
	UnknownNotificationType ChatType = iota
	PersonalChat
	GroupChat
)

const (
	UnknownMessageType MessageType = iota
	StandardMessage
	MatchInvitation
)

const (
	UnknownMessageStatus MessageStatus = iota
	SendingMessageStatus
	FailedMessageStatus
	SentMessageStatus
	DeliveredMessageStatus
	SeenMessageStatus
)
