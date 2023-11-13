package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	appName = "Go Sample"
	appVersion = "v0.2.0"
)

type (
	Config struct {
		AppName string
		AppVersion string
		AppInfo string
		AppPath string `yaml:"app_path"`
		ConfigPath string
		Debug bool `yaml:"debug" env:"APPX_DEBUG"`
		Server `yaml:"server"`
		Listen `yaml:"listen"`
		Log `yaml:"logger"`
		Storage `yaml:"storage"`
		Redis `yaml:"redis"`
		DownloadImages `yaml:"download_images"`
		FileStorage `yaml:"file_storage"`
		S3 `yaml:"s3"`
		// Sentry `yaml:"sentry"`
		Cors		   `yaml:"cors"`
		Translation	`yaml:"translation"`
		ClientSections `yaml:"client_sections"`
		ModulesAccess   `yaml:"modules_access"`
	}

	Server struct {
		ReadTimeout int32 `yaml:"read_timeout"`
		WriteTimeout int32 `yaml:"write_timeout"`
		ShutdownTimeout int32 `yaml:"shutdown_timeout"`
	}

	Listen struct {
		Type string `yaml:"type" env:"APPX_SERVICE_LISTEN_TYPE"`
		SockName string `yaml:"sock_name" env:"APPX_SERVICE_LISTEN_SOCK"`
		BindIP string `yaml:"bind_ip" env:"APPX_SERVICE_BIND"`
		Port string `yaml:"port" env:"APPX_SERVICE_PORT"`
	}

	Log struct {
		Prefix string `yaml:"prefix"`
		Level string `yaml:"level" env:"APPX_LOG_LEVEL"`
	}

	Storage struct {
		Host		string `yaml:"host" env:"APPX_DB_HOST"`
		Port		string `yaml:"port" env:"APPX_DB_PORT"`
		Username	string `yaml:"username" env:"APPX_DB_USER"`
		Password	string `yaml:"password" env:"APPX_DB_PASSWORD"`
		Database	string `yaml:"database" env:"APPX_DB_NAME"`
		MaxPoolSize int32  `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
		Timeout	 int32  `yaml:"timeout"` // in sec
	}

	Redis struct {
		Host string `yaml:"host" env:"APPX_REDIS_HOST"`
		Port string `yaml:"port" env:"APPX_REDIS_PORT"`
		Password string `yaml:"password" env:"APPX_REDIS_PASSWORD"`
		Timeout int `yaml:"timeout"` // in sec
	}

	DownloadImages struct {
		Host string `yaml:"host" env:"APPX_IMAGE_HOST"`
		BasePath string `yaml:"base_path"`
	}

	FileStorage struct {
		DownloadDir string `yaml:"download_dir" env:"APPX_FS_DOWNLOAD_DIR"`
		CatalogCategoryImageDir string `yaml:"catalog_category_image_dir"`
	}

	S3 struct {
		Host string `yaml:"host" env:"APPX_S3_HOST"`
		Port string `yaml:"port" env:"APPX_S3_PORT"`
		UseSSL bool `yaml:"use_ssl" env:"APPX_S3_USESSL"`
		Username string `yaml:"username" env:"APPX_S3_USER"`
		Password string `yaml:"password" env:"APPX_S3_PASSWORD"`
		BucketName string `yaml:"bucket_name"`
	}

	//Sentry struct {
	//	Use bool `yaml:"use" env:"APPX_SENTRY_USE"`
	//	Dsn string `yaml:"dsn" env:"APPX_SENTRY_DSN"`
	//}

	Cors struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
		ExposedHeaders []string `yaml:"exposed_headers"`
		AllowCredentials bool `yaml:"allow_credentials"`
	}

	Translation struct {
		DirPath string `yaml:"dir_path"`
		FileType string `yaml:"file_type"`
		LangCodes []string `yaml:"lang_codes"`
	}

	ClientSections struct {
		AdminAPI ClientSection `yaml:"admin_api"`
		PublicAPI ClientSection `yaml:"public_api"`
	}

	ClientSection struct {
		Caption string `yaml:"caption"`
		Privilege string `yaml:"privilege"`
	}

	ModulesAccess struct {
		Roles `yaml:"roles"`
		Privileges []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}

	Roles struct {
		DirPath string `yaml:"dir_path"`
		FileType string `yaml:"file_type"`
		List []string `yaml:"list"`
	}
)

func New(filePath string) (*Config, error) {
	cfg := &Config{
		AppName: appName,
		AppVersion: appVersion,
		ConfigPath: filePath,
	}

	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		return nil, fmt.Errorf("while reading config '%s', error '%s' occurred", filePath, err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	if cfg.AppPath == "" {
		cfg.AppPath = os.Args[0]
	}

	return cfg, nil
}
