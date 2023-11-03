package http

type StatusCode = uint

const (
	HttpStatusCodeContinue                    = 100
	HttpStatusCodeSwitchingProtocols          = 101
	HttpStatusCodeOK                          = 200
	HttpStatusCodeCreated                     = 201
	HttpStatusCodeAccepted                    = 202
	HttpStatusCodeNonAuthoritativeInformation = 203
	HttpStatusCodeNoContent                   = 204
	HttpStatusCodeResetContent                = 205
	HttpStatusCodePartialContent              = 206
	HttpStatusCodeMultipleChoice              = 300
	HttpStatusCodeMovedPermanently            = 301
	HttpStatusCodeFound                       = 302
	HttpStatusCodeSeeOther                    = 303
	HttpStatusCodeNotModified                 = 304
	HttpStatusCodeUseProxy                    = 305
	HttpStatusCodeTemporaryRedirect           = 307
	HttpStatusCodeBadRequest                  = 400
	HttpStatusCodeUnauthorized                = 401
	HttpStatusCodePaymentRequired             = 402
	HttpStatusCodeForbidden                   = 403
	HttpStatusCodeNotFound                    = 404
	HttpStatusCodeMethodNotAllowed            = 405
	HttpStatusCodeNotAcceptable               = 406
	HttpStatusCodeProxyAuthenticationRequired = 407
	HttpStatusCodeRequestTimeout              = 408
	HttpStatusCodeConflict                    = 409
	HttpStatusCodeGone                        = 410
	HttpStatusCodeLengthRequired              = 411
	HttpStatusCodePreconditionFailed          = 412
	HttpStatusCodePayloadTooLarge             = 413
	HttpStatusCodeURITooLong                  = 414
	HttpStatusCodeUnsupportedMediaType        = 415
	HttpStatusCodeRangeNotSatisfiable         = 416
	HttpStatusCodeExpecationFailed            = 417
	HttpStatusCodeUpgradeRequired             = 426
	HttpStatusCodeInternalServerError         = 500
	HttpStatusCodeNotImplemented              = 501
	HttpStatusCodeBadGateway                  = 502
	HttpStatusCodeServiceUnavailable          = 503
	HttpStatusCodeGatewayTimeout              = 504
	HttpStatusCodeHttpVersionNotSupported     = 505
)
