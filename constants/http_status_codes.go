package constants

const (
	/*
		成功時のステータスコード
	*/
	StatusOK       = 200 // OK
	StatusCreated  = 201 // Created
	StatusAccepted = 202 // Accepted

	/*
		リダイレクト ステータスコード
	*/
	StatusFound = 302 // Status Found

	/*
		クライアントエラー ステータスコード
	*/
	StatusBadRequest       = 400 // Bad Request
	StatusUnauthorized     = 401 // Unauthorized
	StatusForbidden        = 403 // Forbidden
	StatusNotFound         = 404 // Not Found
	StatusMethodNotAllowed = 405 // Method Not Allowed
	StatusConflict         = 409 // Conflict

	/*
		サーバーエラー ステータスコード
	*/
	StatusInternalServerError = 500 // Internal Server Error
	StatusNotImplemented      = 501 // Not Implemented
	StatusBadGateway          = 502 // Bad Gateway
	StatusServiceUnavailable  = 503 // Service Unavailable
	StatusGatewayTimeout      = 504 // Gateway Timeout
)
