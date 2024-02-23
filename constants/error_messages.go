package constants

const (
	InvalidRequest      = "無効なリクエストです"       // 400 Bad Request
	CodeNotFound        = "コードが見つかりません"      // 404 Not Found
	Unauthorized        = "認証に失敗しました"        // 401 Unauthorized
	SecretMismatch      = "シークレットが一致しません"    // 401 Unauthorized
	InternalServerError = "サーバーエラーが発生しました"   // 500 Internal Server Error
	DatabaseError       = "データベースエラーが発生しました" // 500 Internal Server Error
	UnknownError        = "不明なエラーが発生しました"    // 500 Internal Server Error
	Success             = "成功"               // 200 OK
	GroupCodeExists     = "グループコードが存在します"    // 200 OK
	SecretExists        = "シークレットが存在します"     // 200 OK
	GroupCodeVerified   = "グループコードが検証されました"  // 200 OK
)
