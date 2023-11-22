package inconst

const (
	SVC_ACCOUNT  = "account"
	SVC_AUTH     = "auth"
	SVC_PAYMENT  = "payment"
	SVC_SECURITY = "security"

	SEC_ENCRYPT = "/payload/encrypt"
	SEC_DECRYPT = "/payload/decrypt"

	HeaderXPartnerID     = "X-Partner-Id"
	HeaderXSecKeypair    = "X-Sec-Keypair"
	HeaderXCorrelationID = "X-Correlation-Id"
	HeaderXExternalID    = "X-External-Id"
	HeaderXRequestID     = "X-Request-Id"
	HeaderXOverrideSec   = "X-Override-Sec"
)
