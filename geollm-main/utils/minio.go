package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"
	"time"
)
import "github.com/minio/minio-go/v7"

var MinioClient *minio.Client
var err error

func InitMinIO() {
	// Initialize minio client object.
	MinioClient, err = minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//log.Printf("%#v\n", MinioClient) // minioClient is now setup
	CreateBucket()
}

func CreateBucket() {
	err = MinioClient.MakeBucket(context.Background(), "answersheet", minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully created mybucket.")
}

func DeleteFile(aid int, t int) error {
	for i := 0; i <= t; i++ {
		filename := strconv.Itoa(aid) + "_" + strconv.Itoa(t)
		err = MinioClient.RemoveObject(context.Background(), "answersheet", filename, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	for i := 0; i <= t; i++ {
		filename := strconv.Itoa(aid) + "_" + strconv.Itoa(t) + "_thumbnail"
		err = MinioClient.RemoveObject(context.Background(), "answersheet", filename, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("Successfully deleted object")
	return nil
}

func GetFileUrl(aid int, picType int, examType string) ([]string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	//reqParams.Set("response-content-type", "image/jpeg")
	var urls []string
	t, _ := strconv.Atoi(examType)
	for i := 0; i < t+1; i++ {
		objectName := strconv.Itoa(aid) + "_" + strconv.Itoa(i)

		if picType == 1 {
			objectName = objectName + "_thumbnail"
		}
		presignedURL, _ := MinioClient.PresignedGetObject(context.Background(), "answersheet", objectName, time.Second*24*60*60, reqParams)
		urls = append(urls, presignedURL.String())
	}
	// Generates a presigned url which expires in a day.

	return urls, nil
}

func UploadFiles(fileHeaders []*multipart.FileHeader, aid int) error {
	for i, fileHeader := range fileHeaders {
		// 使用文件的基础名称作为对象名称
		objectName := strconv.Itoa(aid) + "_" + strconv.Itoa(i)
		originalFileName := fileHeader.Filename

		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("无法打开文件 %s: %v", originalFileName, err)
		}
		defer file.Close()

		// 将文件内容直接从内存中压缩为原图
		compressedOriginalData, err := CompressImageFromReader(file, RawQuality, 0)
		if err != nil {
			return fmt.Errorf("无法压缩原图 %s: %v", originalFileName, err)
		}

		// 重置文件指针
		file.Seek(0, io.SeekStart)

		// 压缩缩略图
		compressedThumbnailData, err := CompressImageFromReader(file, ThumbnailQuality, 1)
		if err != nil {
			return fmt.Errorf("无法压缩缩略图 %s: %v", originalFileName, err)
		}

		// 上传原图
		err = uploadToMinioFromBytes(compressedOriginalData, objectName, "answersheet")
		if err != nil {
			return fmt.Errorf("上传原图失败: %v", err)
		}

		// 上传缩略图
		err = uploadToMinioFromBytes(compressedThumbnailData, objectName+"_thumbnail", "answersheet")
		if err != nil {
			return fmt.Errorf("上传缩略图失败: %v", err)
		}
	}
	return nil
}

// CompressImageFromReader 压缩图片并调整大小以控制文件大小
func CompressImageFromReader(file io.Reader, quality, re_size int) ([]byte, error) {
	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("无法解码图片: %v", err)
	}

	if re_size == 1 {
		// 获取图片的原始尺寸
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()

		// 缩放到宽度 800 或更小，以减少文件大小
		maxWidth := 800
		if width > maxWidth {
			scale := float64(maxWidth) / float64(width)
			newWidth := uint(float64(width) * scale)
			newHeight := uint(float64(height) * scale)
			img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
		}
	}
	// 压缩图片
	compressedImage := new(bytes.Buffer)
	err = jpeg.Encode(compressedImage, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("无法压缩图片: %v", err)
	}

	return compressedImage.Bytes(), nil
}

// 上传字节数据到 Minio 的辅助函数
func uploadToMinioFromBytes(data []byte, objectName, bucketName string) error {
	// 根据对象名称获取 MIME 类型
	contentType := "image/jpeg"

	// 创建一个临时文件来模拟上传对象
	uploadInfo, err := MinioClient.PutObject(
		context.Background(),
		bucketName, // 存储桶名称
		objectName, // 对象名称
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("无法上传文件到 Minio: %v", err)
	}

	fmt.Println("Successfully uploaded object:", uploadInfo)
	return nil
}
