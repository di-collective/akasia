package config

type Environment struct {
	DbSslMode          string `env:"DB_SSL_MODE" envDefault:"disable"`
	DbHost             string `env:"DB_HOST" envDefault:"localhost"`
	DbName             string `env:"DB_NAME" envDefault:"postgres"`
	DbPort             string `env:"DB_PORT" envDefault:"5432"`
	DbUser             string `env:"DB_USER" envDefault:"postgres"`
	DbPass             string `env:"DB_PASS" envDefault:"sa123456"`
	FirebaseConfig     string `env:"FIREBASE_CONFIG" envDefault:"./firebase.json"`
	JWTAlgo            string `env:"JWT_ALGO" envDefault:"HS256"`
	JWTSecret          string `env:"JWT_SECRET" envDefault:"JWTSECRET11234567892123456789312"`
	ServicePort        int    `env:"SVC_PORT" envDefault:"3333"`
	SMTPHost           string `env:"SMTP_HOST" envDefault:"smtpdm-ap-southeast-1.aliyun.com"`
	SMTPPort           int    `env:"SMTP_PORT" envDefault:"465"`
	SMTPAuthEmail      string `env:"SMTP_AUTH_EMAIL" envDefault:"no-reply@edufund.co.id"`
	SMTPAuthPassword   string `env:"SMTP_AUTH_PASSWORD" envDefault:"5PHT57g5YUxqgPHMa"`
	CsMail             string `env:"CS_MAIL" envDefault:"customer.care@edufund.co.id"`
	ResetPasswordUrl   string `env:"RESET_PASSWORD_URL" envDefault:"https://edufund.co.id"`
	DirPath            string `env:"DIR_PATH" envDefault:"/User/template"`
	OSSEndpoint        string `env:"OSS_ENDPOINT" envDefault:"oss-ap-southeast.aliyuncs.com"`
	OSSAccessKeyID     string `env:"OSS_ACCESS_KEY_ID" envDefault:"12345678901234567890ABCD"`
	OSSAccessKeySecret string `env:"OSS_ACCESS_KEY_SECRET" envDefault:"1234567890ABCDE1234567890ABCDE"`
	OSSBucketName      string `env:"OSS_BUCKET_NAME" envDefault:"bucket-name"`
}
