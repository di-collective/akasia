package config

type Environment struct {
	DbSslMode          string `env:"DB_SSL_MODE"`
	DbHost             string `env:"DB_HOST"`
	DbName             string `env:"DB_NAME"`
	DbPort             string `env:"DB_PORT"`
	DbUser             string `env:"DB_USER"`
	DbPass             string `env:"DB_PASS"`
	FirebaseConfig     string `env:"FIREBASE_CONFIG"`
	JWTAlgo            string `env:"JWT_ALGO"`
	JWTSecret          string `env:"JWT_SECRET"`
	ServicePort        int    `env:"SVC_PORT"`
	SMTPHost           string `env:"SMTP_HOST"`
	SMTPPort           int    `env:"SMTP_PORT"`
	SMTPAuthEmail      string `env:"SMTP_AUTH_EMAIL"`
	SMTPAuthPassword   string `env:"SMTP_AUTH_PASSWORD"`
	CsMail             string `env:"CS_MAIL"`
	ResetPasswordUrl   string `env:"RESET_PASSWORD_URL"`
	DirPath            string `env:"DIR_PATH"`
	OSSEndpoint        string `env:"OSS_ENDPOINT"`
	OSSAccessKeyID     string `env:"OSS_ACCESS_KEY_ID"`
	OSSAccessKeySecret string `env:"OSS_ACCESS_KEY_SECRET"`
	OSSBucketName      string `env:"OSS_BUCKET_NAME"`
	BaseURLUser        string `env:"BASE_URL_USER"`
	Capacity           string `env:"CAPACITY"`
	BaseURLClinic      string `env:"BASE_URL_CLINIC"`
}
