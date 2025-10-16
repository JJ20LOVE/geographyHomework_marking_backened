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

func init() {
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

	// 从环境变量读取OCR配置
	OcrProvider = os.Getenv("OCR_PROVIDER")
	if OcrProvider == "" {
		OcrProvider = "aliyun" // 默认使用阿里云
	}

	AliyunAccessKeyId = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	AliyunAccessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	AliyunRegion = os.Getenv("ALIYUN_REGION")
	if AliyunRegion == "" {
		AliyunRegion = "cn-hangzhou" // 默认区域
	}

	// 这里也可以从配置文件读取，根据你的配置加载方式
	fmt.Printf("OCR配置: Provider=%s, Region=%s\n", OcrProvider, AliyunRegion)
	EvaApi = viper.GetString("eva.EvaApi")

	Endpoint = viper.GetString("minio.Endpoint")
	AccessKeyID = viper.GetString("minio.AccessKeyID")
	SecretAccessKey = viper.GetString("minio.SecretAccessKey")
	RawQuality = viper.GetInt("minio.RawQuality")
	ThumbnailQuality = viper.GetInt("minio.ThumbnailQuality")

	// AI配置 - 新增
	AIProvider = viper.GetString("ai.provider")
	if AIProvider == "" {
		AIProvider = "deepseek" // 默认使用deepseek
	}

	DeepSeekApiKey = viper.GetString("ai.deepseek.api_key")
	DeepSeekBaseUrl = viper.GetString("ai.deepseek.base_url")
	if DeepSeekBaseUrl == "" {
		DeepSeekBaseUrl = "https://api.deepseek.com" // 默认URL
	}

	AIModel = viper.GetString("ai.model")
	if AIModel == "" {
		AIModel = "deepseek-chat" // 默认模型
	}

	AIMaxTokens = viper.GetInt("ai.max_tokens")
	if AIMaxTokens == 0 {
		AIMaxTokens = 2000 // 默认最大token数
	}

	AITemperature = viper.GetFloat64("ai.temperature")
	if AITemperature == 0 {
		AITemperature = 0.7 // 默认温度参数
	}

	fmt.Printf("AI配置: Provider=%s, Model=%s\n", AIProvider, AIModel)
}
