package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	AppMode          string
	HttpPort         string
	JwtKey           string
	Db               string
	DbHost           string
	DbPort           string
	DbUser           string
	DbPassWord       string
	DbName           string
	MaxOpenConns     int
	MaxIdleConns     int
	OcrApi           string
	FilePrefix       string
	EvaApi           string
	Endpoint         string
	AccessKeyID      string
	SecretAccessKey  string
	RawQuality       int
	ThumbnailQuality int
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

	OcrApi = viper.GetString("ocr.OcrApi")
	EvaApi = viper.GetString("eva.EvaApi")

	Endpoint = viper.GetString("minio.Endpoint")
	AccessKeyID = viper.GetString("minio.AccessKeyID")
	SecretAccessKey = viper.GetString("minio.SecretAccessKey")
	RawQuality = viper.GetInt("minio.RawQuality")
	ThumbnailQuality = viper.GetInt("minio.ThumbnailQuality")
}
