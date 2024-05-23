package constants

const (
	DEV  = "development"
	HOST = "localhost"
)

const (
	USD = "USD"
	EUR = "EUR"
	VND = "VND"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

const (
	ErrInvalidToken = "Token Is Invalid"
	ErrExpiredToken = "Token Has Expired"
)

const (
	KeyRefetchToken = "refetchToken"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

const (
	DepositorRole = "depositor"
	BankerRole    = "banker"
)
