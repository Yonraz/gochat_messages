package constants

type Queues string
type RoutingKey string
type Exchange string
type UserStatus string
type Notification string
type MessageType string

const (
	UserRegisteredKey RoutingKey = "user.registered"
)

const (
	UserRegistrationQueue Queues = "UserRegistrationQueue"
)

const (
	Online  UserStatus = "online"
	Offline UserStatus = "offline"
)

const (
	MessageUpdate MessageType = "message.update"
	MessageCreate MessageType = "message.create"
)

const (
	UserEventsExchange    Exchange = "UserEventsExchange"
	MessageEventsExchange Exchange = "MessageEventsExchange"
)

const (
	UserLoggedInKey     RoutingKey = "user.logged.in"
	UserSignedoutKey    RoutingKey = "user.signed.out"
	MessageSentKey      RoutingKey = "message.sent"
	MessageDeliveredKey RoutingKey = "message.delivered"
	MessageReadKey      RoutingKey = "message.read"
)

const (
	UserLoginQueue        Queues = "MESSAGES_SRV_UserLoginQueue"
	UserSignoutQueue      Queues = "MESSAGES_SRV_UserSignoutQueue"
	MessageSentQueue      Queues = "MESSAGES_SRV_MessageSentQueue"
	MessageDeliveredQueue Queues = "MESSAGES_SRV_MessageDeliveredQueue"
	MessageReadQueue      Queues = "MESSAGES_SRV_MessageReadQueue"
)

const (
	UserSentMessage Notification = "user.sent.message"
)
