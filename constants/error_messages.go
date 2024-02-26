package constants

const (
	InvalidRequest           = "無効なリクエストです"       // 400 Bad Request
	BadRequestMessage        = "リクエストが不正です"       // 400 Bad Request
	CodeNotFound             = "コードが見つかりません"      // 404 Not Found
	Unauthorized             = "認証に失敗しました"        // 401 Unauthorized
	SecretMismatch           = "シークレットが一致しません"    // 401 Unauthorized
	InternalServerError      = "サーバーエラーが発生しました"   // 500 Internal Server Error
	DatabaseError            = "データベースエラーが発生しました" // 500 Internal Server Error
	UnknownError             = "不明なエラーが発生しました"    // 500 Internal Server Error
	Success                  = "成功"               // 200 OK
	GroupCodeExists          = "グループコードが存在します"    // 200 OK
	SecretExists             = "シークレットが存在します"     // 200 OK
	GroupCodeVerified        = "グループコードが検証されました"  // 200 OK
	ErrNoFileHeaderJP        = "ファイルヘッダが提供されていません"
	ErrOpenFileJP            = "ファイルのオープンに失敗しました"
	ErrReadFileDataJP        = "ファイルデータの読み取りに失敗しました"
	ErrLoadAWSConfigJP       = "AWS設定のロードに失敗しました"
	ErrUploadToS3JP          = "S3へのアップロードに失敗しました"
	ErrCloudFrontURLNotSetJP = "AWS_CLOUDFRONT環境変数が設定されていません"
)
