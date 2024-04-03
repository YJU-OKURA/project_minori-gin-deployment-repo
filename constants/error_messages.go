package constants

// クライアントエラー関連のエラーメッセージ
const (
	InvalidRequest    = "無効なリクエストです"          // 400 Bad Request
	BadRequestMessage = "リクエストが不正です"          // 400 Bad Request
	ErrNoFileHeaderJP = "ファイルヘッダが提供されていません"   // 400 Bad Request
	ErrFileSizeJP     = "ファイルサイズが10MBを超えています" // 400 Bad Request
	ErrMimeTypeJP     = "ファイルタイプが画像ではありません"   // 400 Bad Request
	ErrNoDateJP       = "日付が提供されていません"        // 400 Bad Request
	ErrInvalidInput   = "無効な入力です"             // 400 Bad Request
	ErrNoUserID       = "ユーザーIDが提供されていません"    // 400 Bad Request
)

// 認証関連のエラーメッセージ
const (
	Unauthorized          = "認証に失敗しました"       // 401 Unauthorized
	SecretMismatch        = "シークレットが一致しません"   // 401 Unauthorized
	CodeNotFound          = "コードが見つかりません"     // 404 Not Found
	ClassNotFound         = "クラスが見つかりません"     // 404 Not Found
	ApplyingClassNotFound = "申請中のクラスが見つかりません" // 404 Not Found
	UserNotFound          = "ユーザーが見つかりません"    // 404 Not Found
)

// サーバーエラー&データベース関連のエラーメッセージ
const (
	InternalServerError      = "サーバーエラーが発生しました"               // 500 Internal Server Error
	DatabaseError            = "データベースエラーが発生しました"             // 500 Internal Server Error
	UnknownError             = "不明なエラーが発生しました"                // 500 Internal Server Error
	ErrOpenFileJP            = "ファイルのオープンに失敗しました"             // 500 Internal Server Error
	ErrReadFileDataJP        = "ファイルデータの読み取りに失敗しました"          // 500 Internal Server Error
	ErrLoadAWSConfigJP       = "AWS設定のロードに失敗しました"             // 500 Internal Server Error
	ErrUploadToS3JP          = "S3へのアップロードに失敗しました"            // 500 Internal Server Error
	ErrCloudFrontURLNotSetJP = "AWS_CLOUDFRONT環境変数が設定されていません" // 500 Internal Server Error
	AssignError              = "ロールの割り当てに失敗しました"              // 500 Internal Server Error
	ErrLoadMessage           = "メッセージの取得に失敗しました"              // 500 Internal Server Error
	ErrSendMessage           = "メッセージの送信に失敗しました"              // 500 Internal Server Error
)

// 成功時のメッセージ
const (
	Success                 = "成功"                // 200 OK
	ClassCodeExists         = "クラスコードが存在します"      // 200 OK
	SecretExists            = "シークレットが存在します"      // 200 OK
	ClassCodeVerified       = "クラスコードが検証されました"    // 200 OK
	ClassMemberRegistration = "クラスコードの確認と役割の割り当て" // 200 OK
	CreateOrUpdateSuccess   = "作成または更新に成功しました"    // 200 OK
	DeleteSuccess           = "削除に成功しました"         // 200 OK
	MessageSent             = "メッセージが送信されました"     // 200 OK
)
