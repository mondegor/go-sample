package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	appName    = "Go Sample"
	appVersion = "v0.5.0"
)

type (
	Config struct {
		AppName       string
		AppVersion    string
		AppInfo       string
		AppPath       string `yaml:"app_path"`
		ConfigPath    string
		Debugging     `yaml:"debugging"`
		Log           `yaml:"logger"`
		Server        `yaml:"server"`
		Listen        `yaml:"listen"`
		Storage       `yaml:"storage"`
		Redis         `yaml:"redis"`
		FileSystem    `yaml:"file_system"`
		S3            `yaml:"s3"`
		FileProviders `yaml:"file_providers"`
		// Sentry        `yaml:"sentry"`
		Cors            `yaml:"cors"`
		Translation     `yaml:"translation"`
		ClientSections  `yaml:"client_sections"`
		ModulesAccess   `yaml:"modules_access"`
		ModulesSettings `yaml:"modules_settings"`
	}

	Debugging struct {
		Debug       bool `yaml:"debug" env:"APPX_DEBUG"`
		ErrorCaller `yaml:"caller"`
	}

	ErrorCaller struct {
		Deep         int    `yaml:"deep" env:"APPX_ERR_CALLER_DEEP"`
		UseShortPath bool   `yaml:"use_short_path" env:"APPX_ERR_CALLER_USE_SHORT_PATH"`
		RootPath     string `yaml:"root_path"`
	}

	Log struct {
		Prefix    string `yaml:"prefix" env:"APPX_LOG_PREFIX"`
		Level     string `yaml:"level" env:"APPX_LOG_LEVEL"`
		LogCaller `yaml:"caller"`
	}

	LogCaller struct {
		Deep         int    `yaml:"deep" env:"APPX_LOG_CALLER_DEEP"`
		UseShortPath bool   `yaml:"use_short_path" env:"APPX_LOG_CALLER_USE_SHORT_PATH"`
		RootPath     string `yaml:"root_path"`
	}

	Server struct {
		ReadTimeout     int32 `yaml:"read_timeout" env:"APPX_SERVER_READ_TIMEOUT"`
		WriteTimeout    int32 `yaml:"write_timeout" env:"APPX_SERVER_WRITE_TIMEOUT"`
		ShutdownTimeout int32 `yaml:"shutdown_timeout" env:"APPX_SERVER_SHUTDOWN_TIMEOUT"`
	}

	Listen struct {
		Type     string `yaml:"type" env:"APPX_SERVICE_LISTEN_TYPE"`
		SockName string `yaml:"sock_name" env:"APPX_SERVICE_LISTEN_SOCK"`
		BindIP   string `yaml:"bind_ip" env:"APPX_SERVICE_BIND"`
		Port     string `yaml:"port" env:"APPX_SERVICE_PORT"`
	}

	Storage struct {
		Host        string `yaml:"host" env:"APPX_DB_HOST"`
		Port        string `yaml:"port" env:"APPX_DB_PORT"`
		Username    string `yaml:"username" env:"APPX_DB_USER"`
		Password    string `yaml:"password" env:"APPX_DB_PASSWORD"`
		Database    string `yaml:"database" env:"APPX_DB_NAME"`
		MaxPoolSize int32  `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
		Timeout     int32  `yaml:"timeout"` // in sec
	}

	Redis struct {
		Host     string `yaml:"host" env:"APPX_REDIS_HOST"`
		Port     string `yaml:"port" env:"APPX_REDIS_PORT"`
		Password string `yaml:"password" env:"APPX_REDIS_PASSWORD"`
		Timeout  int    `yaml:"timeout"` // in sec
	}

	FileSystem struct {
		DirMode    uint32 `yaml:"dir_mode" env:"APPX_FILESYSTEM_DIR_MODE"`
		CreateDirs bool   `yaml:"create_dirs" env:"APPX_FILESYSTEM_CREATE_DIRS"`
	}

	S3 struct {
		Host          string `yaml:"host" env:"APPX_S3_HOST"`
		Port          string `yaml:"port" env:"APPX_S3_PORT"`
		UseSSL        bool   `yaml:"use_ssl" env:"APPX_S3_USESSL"`
		Username      string `yaml:"username" env:"APPX_S3_USER"`
		Password      string `yaml:"password" env:"APPX_S3_PASSWORD"`
		CreateBuckets bool   `yaml:"create_buckets" env:"APPX_S3_CREATE_BUCKETS"`
	}

	FileProviders struct {
		ImageStorage struct {
			Name       string `yaml:"name"`
			BucketName string `yaml:"bucket_name" env:"APPX_IMAGESTORAGE_BUCKET"`
		} `yaml:"image_storage"`
		ImageStorage2 struct {
			Name    string `yaml:"name"`
			RootDir string `yaml:"root_dir" env:"APPX_IMAGESTORAGE2_ROOT_DIR"`
		} `yaml:"image_storage2"`
	}

	//Sentry struct {
	//	Use bool `yaml:"use" env:"APPX_SENTRY_USE"`
	//	Dsn string `yaml:"dsn" env:"APPX_SENTRY_DSN"`
	//}

	Cors struct {
		AllowedOrigins   []string `yaml:"allowed_origins" env:"APPX_CORS_ALLOWED_ORIGINS"` // items by "," separated
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
	}

	Translation struct {
		DirPath   string   `yaml:"dir_path"`
		FileType  string   `yaml:"file_type"`
		LangCodes []string `yaml:"lang_codes" env:"APPX_TRANSLATION_LANGS"` // items by "," separated
	}

	ClientSections struct {
		AdminAPI  ClientSection `yaml:"admin_api"`
		PublicAPI ClientSection `yaml:"public_api"`
	}

	ClientSection struct {
		Caption   string `yaml:"caption"`
		Privilege string `yaml:"privilege"`
	}

	ModulesAccess struct {
		Roles       `yaml:"roles"`
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}

	Roles struct {
		DirPath  string   `yaml:"dir_path"`
		FileType string   `yaml:"file_type"`
		List     []string `yaml:"list"`
	}

	ModulesSettings struct {
		CatalogCategory struct {
			Image struct {
				BaseDir      string `yaml:"base_dir"`
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage or ImageStorage2
			} `yaml:"image"`
		} `yaml:"catalog_category"`
		FileStation struct {
			ImageProxy struct {
				Host         string `yaml:"host" env:"APPX_IMAGE_HOST"`
				BaseURL      string `yaml:"base_url"`
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage or ImageStorage2
			} `yaml:"image_proxy"`
		} `yaml:"file_station"`
	}
)

func New(filePath string) (*Config, error) {
	cfg := &Config{
		AppName:    appName,
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
