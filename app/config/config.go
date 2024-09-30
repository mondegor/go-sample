package config

import (
	"io"
	"time"

	"github.com/mondegor/go-webcore/mrlib"
)

type (
	// Args - разобранные аргументы, которые передаются из командной строки.
	// Эти аргументы более приоритетны аналогичным, которые определены в конфигурации или переданы через .env файл.
	Args struct {
		WorkDir     string // путь к рабочей директории приложения
		ConfigPath  string // путь к файлу конфигурации приложения
		DotEnvPath  string // путь к .env файлу (переменные из этого файла более приоритетны переменных из ConfigPath)
		Environment string // внешнее окружение: local, dev, test, prod
		LogLevel    string // уровень логирования: info, warn, error, fatal, debug, trace
		Stdout      io.Writer
	}

	// Config - comment struct.
	Config struct {
		Os
		App             `yaml:"app"`
		Debugging       `yaml:"debugging"`
		Log             `yaml:"logger"`
		Sentry          `yaml:"sentry"`
		Servers         `yaml:"servers"`
		Storage         `yaml:"storage"`
		Redis           `yaml:"redis"`
		FileSystem      `yaml:"file_system"`
		S3              `yaml:"s3"`
		FileProviders   `yaml:"file_providers"`
		Cors            `yaml:"cors"`
		Translation     `yaml:"translation"`
		AppSections     `yaml:"app_sections"`
		AccessControl   `yaml:"access_control"`
		ModulesSettings `yaml:"modules_settings"`
		MimeTypes       `yaml:"mime_types"`
		TaskSchedule    `yaml:"task_schedule"`
	}

	// Os - comment struct.
	Os struct {
		Stdout io.Writer
	}

	// App - comment struct.
	App struct {
		Name        string `yaml:"name" env:"APPX_NAME"`
		Version     string `yaml:"version" env:"APPX_VER"`
		Environment string `yaml:"environment" env:"APPX_ENV"`
		WorkDir     string
		ConfigPath  string
		DotEnvPath  string
		StartedAt   time.Time
	}

	// Debugging - comment struct.
	Debugging struct {
		Debug                bool `yaml:"debug" env:"APPX_DEBUG"`
		UnexpectedHttpStatus int  `yaml:"unexpected_http_status"`
		ErrorCaller          `yaml:"error_caller"`
	}

	// ErrorCaller - comment struct.
	ErrorCaller struct {
		Enable       bool     `yaml:"enable" env:"APPX_ERR_CALLER_ENABLE"`
		Depth        int      `yaml:"depth" env:"APPX_ERR_CALLER_DEPTH"`
		ShowFuncName bool     `yaml:"show_func_name"`
		UpperBounds  []string `yaml:"upper_bounds"`
	}

	// Log - comment struct.
	Log struct {
		Level           string `yaml:"level" env:"APPX_LOG_LEVEL"`
		TimestampFormat string `yaml:"timestamp_format" env:"APPX_LOG_TIMESTAMP"`
		JsonFormat      bool   `yaml:"json_format" env:"APPX_LOG_JSON"`
		ConsoleColor    bool   `yaml:"console_color" env:"APPX_LOG_COLOR"`
	}

	// Sentry - comment struct.
	Sentry struct {
		DSN              string        `yaml:"dsn" env:"APPX_SENTRY_DSN"`
		TracesSampleRate float64       `yaml:"traces_sample_rate" env:"APPX_SENTRY_TRACES_SAMPLE_RATE"`
		FlushTimeout     time.Duration `yaml:"flush_timeout"`
	}

	// Servers - comment struct.
	Servers struct {
		// RestServer - comment struct.
		RestServer struct {
			ReadTimeout     time.Duration `yaml:"read_timeout" env:"APPX_SERVER_READ_TIMEOUT"`
			WriteTimeout    time.Duration `yaml:"write_timeout" env:"APPX_SERVER_WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APPX_SERVER_SHUTDOWN_TIMEOUT"`
			Listen          struct {
				BindIP string `yaml:"bind_ip" env:"APPX_SERVER_LISTEN_BIND"`
				Port   string `yaml:"port" env:"APPX_SERVER_LISTEN_PORT"`
			} `yaml:"listen"`
		} `yaml:"rest_server"`

		// InternalServer - comment struct.
		InternalServer struct {
			ReadTimeout     time.Duration `yaml:"read_timeout" env:"APPX_INTERNAL_SERVER_READ_TIMEOUT"`
			WriteTimeout    time.Duration `yaml:"write_timeout" env:"APPX_INTERNAL_SERVER_WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APPX_INTERNAL_SERVER_SHUTDOWN_TIMEOUT"`
			Listen          struct {
				BindIP string `yaml:"bind_ip" env:"APPX_INTERNAL_SERVER_LISTEN_BIND"`
				Port   string `yaml:"port" env:"APPX_INTERNAL_SERVER_LISTEN_PORT"`
			} `yaml:"listen"`
		} `yaml:"internal_server"`
	}

	// Storage - comment struct.
	Storage struct {
		Host            string        `yaml:"host" env:"APPX_DB_HOST"`
		Port            string        `yaml:"port" env:"APPX_DB_PORT"`
		Username        string        `yaml:"username" env:"APPX_DB_USER"`
		Password        string        `yaml:"password" env:"APPX_DB_PASSWORD"`
		Database        string        `yaml:"database" env:"APPX_DB_NAME"`
		MigrationsDir   string        `yaml:"migrations_dir"`
		MigrationsTable string        `yaml:"migrations_table" env:"APPX_DB_MIGRATIONS_TABLE"`
		MaxPoolSize     int           `yaml:"max_pool_size" env:"APPX_DB_MAX_POOL_SIZE"`
		MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"APPX_DB_MAX_CONN_LIFETIME"`
		MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" env:"APPX_DB_MAX_CONN_IDLE_TIME"`
		Timeout         time.Duration `yaml:"timeout"`
	}

	// Redis - comment struct.
	Redis struct {
		Host         string        `yaml:"host" env:"APPX_REDIS_HOST"`
		Port         string        `yaml:"port" env:"APPX_REDIS_PORT"`
		Password     string        `yaml:"password" env:"APPX_REDIS_PASSWORD"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"APPX_REDIS_READ_TIMEOUT"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"APPX_REDIS_WRITE_TIMEOUT"`
	}

	// FileSystem - comment struct.
	FileSystem struct {
		DirMode    uint32 `yaml:"dir_mode" env:"APPX_FILESYSTEM_DIR_MODE"`
		CreateDirs bool   `yaml:"create_dirs" env:"APPX_FILESYSTEM_CREATE_DIRS"`
	}

	// S3 - comment struct.
	S3 struct {
		Host          string `yaml:"host" env:"APPX_S3_HOST"`
		Port          string `yaml:"port" env:"APPX_S3_PORT"`
		UseSSL        bool   `yaml:"use_ssl" env:"APPX_S3_USESSL"`
		Username      string `yaml:"username" env:"APPX_S3_USER"`
		Password      string `yaml:"password" env:"APPX_S3_PASSWORD"`
		CreateBuckets bool   `yaml:"create_buckets" env:"APPX_S3_CREATE_BUCKETS"`
	}

	// FileProviders - comment struct.
	FileProviders struct {
		// ImageStorage - comment struct.
		ImageStorage struct {
			Name       string `yaml:"name"`
			BucketName string `yaml:"bucket_name" env:"APPX_IMAGESTORAGE_BUCKET"`
		} `yaml:"image_storage"`
		ImageStorage2 struct {
			Name    string `yaml:"name"`
			RootDir string `yaml:"root_dir" env:"APPX_IMAGESTORAGE2_ROOT_DIR"`
		} `yaml:"image_storage2"`
	}

	// Cors - comment struct.
	Cors struct {
		AllowedOrigins   []string `yaml:"allowed_origins" env:"APPX_CORS_ALLOWED_ORIGINS"` // items by "," separated
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
	}

	// Translation - comment struct.
	Translation struct {
		DirPath   string   `yaml:"dir_path" env:"APPX_TRANSLATION_DIR_PATH"`
		LangCodes []string `yaml:"lang_codes" env:"APPX_TRANSLATION_LANGS"` // items by "," separated
		// Dictionaries - comment struct.
		Dictionaries struct {
			DirPath string   `yaml:"dir_path" env:"APPX_TRANSLATION_DICTIONARIES_DIR_PATH"`
			List    []string `yaml:"list"`
		} `yaml:"dictionaries"`
	}

	// AppSections - comment struct.
	AppSections struct {
		// AdminAPI - comment struct.
		AdminAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_ADMIN_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_ADMIN_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"admin_api"`
		// PublicAPI - comment struct.
		PublicAPI struct {
			Privilege string `yaml:"privilege"`
			Auth      struct {
				Secret   string `yaml:"secret" env:"APPX_PUBLIC_API_AUTH_SECRET"`
				Audience string `yaml:"audience" env:"APPX_PUBLIC_API_AUTH_AUDIENCE"`
			} `yaml:"auth"`
		} `yaml:"public_api"`
	}

	// AccessControl - comment struct.
	AccessControl struct {
		Roles       `yaml:"roles"`
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}

	// Roles - comment struct.
	Roles struct {
		DirPath  string   `yaml:"dir_path" env:"APPX_ROLES_DIR_PATH"`
		FileType string   `yaml:"file_type"`
		List     []string `yaml:"list"`
	}

	// ModulesSettings - comment struct.
	ModulesSettings struct {
		// General - comment struct.
		General struct {
			PageSizeMax     uint64 `yaml:"page_size_max"`
			PageSizeDefault uint64 `yaml:"page_size_default"`
		} `yaml:"general"`
		// CatalogCategory - comment struct.
		CatalogCategory struct {
			// Image - comment struct.
			Image struct {
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage or ImageStorage2
			} `yaml:"image"`
		} `yaml:"catalog_category"`
		// FileStation - comment struct.
		FileStation struct {
			// ImageProxy - comment struct.
			ImageProxy struct {
				Host         string `yaml:"host" env:"APPX_IMAGE_HOST"`
				BasePath     string `yaml:"base_path"`
				FileProvider string `yaml:"file_provider"` // FileProviders.ImageStorage or ImageStorage2
			} `yaml:"image_proxy"`
		} `yaml:"file_station"`
	}

	// TaskSchedule - comment struct.
	TaskSchedule struct {
		SettingsReloader SchedulerTask `yaml:"settings_reloader"`
	}

	// SchedulerTask - comment struct.
	SchedulerTask struct {
		Caption string        `yaml:"caption"`
		Startup bool          `yaml:"startup"`
		Period  time.Duration `yaml:"period"`
		Timeout time.Duration `yaml:"timeout"`
	}

	// MimeTypes - comment struct.
	MimeTypes []mrlib.MimeType
)
