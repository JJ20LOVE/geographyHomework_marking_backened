package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	AppMode               string
	HttpPort              string
	JwtKey                string
	Db                    string
	DbHost                string
	DbPort                string
	DbUser                string
	DbPassWord            string
	DbName                string
	MaxOpenConns          int
	MaxIdleConns          int
	OcrApi                string
	FilePrefix            string
	EvaApi                string
	Endpoint              string
	AccessKeyID           string
	SecretAccessKey       string
	RawQuality            int
	ThumbnailQuality      int
	OcrProvider           string
	AliyunAccessKeyId     string
	AliyunAccessKeySecret string
	AliyunRegion          string
	OcrTimeout            int
	AIProvider            string
	DeepSeekApiKey        string
	DeepSeekBaseUrl       string
	AIModel               string
	AIMaxTokens           int
	AITemperature         float64
)

func InitConfig() {
	// 设置默认配置路径
	viper.AddConfigPath("./conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件: %s", err))
	}

	// 设置环境变量读取前缀（可选）
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()

	// 绑定配置项到变量
	AppMode = viper.GetString("server.AppMode")
	HttpPort = viper.GetString("server.HttpPort")
	JwtKey = viper.GetString("server.JwtKey")

	Db = viper.GetString("database.Db")
	DbHost = viper.GetString("database.DbHost")
	DbPort = viper.GetString("database.DbPort")
	DbUser = viper.GetString("database.DbUser")
	DbPassWord = viper.GetString("database.DbPassWord")
	DbName = viper.GetString("database.Dbname")
	MaxOpenConns = viper.GetInt("database.MaxOpenConns")
	MaxIdleConns = viper.GetInt("database.MaxIdleConns")

	// 从环境变量读取OCR配置（现在环境变量已经加载）
	OcrProvider = os.Getenv("OCR_PROVIDER")
	if OcrProvider == "" {
		OcrProvider = viper.GetString("ocr.Provider") // 回退到配置文件
		if OcrProvider == "" {
			OcrProvider = "aliyun" // 最终默认值
		}
	}

	AliyunAccessKeyId = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	if AliyunAccessKeyId == "" {
		AliyunAccessKeyId = viper.GetString("ocr.AliyunAccessKeyId") // 回退到配置文件
	}

	AliyunAccessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	if AliyunAccessKeySecret == "" {
		AliyunAccessKeySecret = viper.GetString("ocr.AliyunAccessKeySecret") // 回退到配置文件
	}

	AliyunRegion = os.Getenv("ALIYUN_REGION")
	if AliyunRegion == "" {
		AliyunRegion = viper.GetString("ocr.AliyunRegion") // 回退到配置文件
		if AliyunRegion == "" {
			AliyunRegion = "cn-hangzhou" // 默认区域
		}
	}

	// 添加调试信息
	fmt.Printf("OCR配置调试信息:\n")
	fmt.Printf("  OcrProvider: %s\n", OcrProvider)
	fmt.Printf("  AliyunAccessKeyId: %s (长度: %d)\n", AliyunAccessKeyId, len(AliyunAccessKeyId))
	fmt.Printf("  AliyunAccessKeySecret: %s (长度: %d)\n", AliyunAccessKeySecret, len(AliyunAccessKeySecret))
	fmt.Printf("  AliyunRegion: %s\n", AliyunRegion)

	EvaApi = viper.GetString("eva.EvaApi")

	Endpoint = viper.GetString("minio.Endpoint")
	AccessKeyID = viper.GetString("minio.AccessKeyID")
	SecretAccessKey = viper.GetString("minio.SecretAccessKey")
	RawQuality = viper.GetInt("minio.RawQuality")
	ThumbnailQuality = viper.GetInt("minio.ThumbnailQuality")

	// AI配置
	AIProvider = viper.GetString("ai.provider")
	if AIProvider == "" {
		AIProvider = "deepseek"
	}

	DeepSeekApiKey = viper.GetString("ai.deepseek.api_key")
	DeepSeekBaseUrl = viper.GetString("ai.deepseek.base_url")
	if DeepSeekBaseUrl == "" {
		DeepSeekBaseUrl = "https://api.deepseek.com"
	}

	AIModel = viper.GetString("ai.model")
	if AIModel == "" {
		AIModel = "deepseek-chat"
	}

	AIMaxTokens = viper.GetInt("ai.max_tokens")
	if AIMaxTokens == 0 {
		AIMaxTokens = 2000
	}

	AITemperature = viper.GetFloat64("ai.temperature")
	if AITemperature == 0 {
		AITemperature = 0.7
	}

	fmt.Printf("AI配置: Provider=%s, Model=%s\n", AIProvider, AIModel)
}
