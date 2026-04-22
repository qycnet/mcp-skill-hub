package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config 存储配置
type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// ObjectStorage 对象存储接口
type ObjectStorage interface {
	// Upload 上传文件
	Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) error
	
	// Download 下载文件
	Download(ctx context.Context, objectName string) (io.Reader, error)
	
	// Delete 删除文件
	Delete(ctx context.Context, objectName string) error
	
	// Exists 检查文件是否存在
	Exists(ctx context.Context, objectName string) (bool, error)
	
	// GetURL 获取下载 URL
	GetURL(ctx context.Context, objectName string) (string, error)
}

// MinIOStorage MinIO 实现
type MinIOStorage struct {
	client *minio.Client
	config Config
}

// NewMinIOStorage 创建 MinIO 存储
func NewMinIOStorage(config Config) (ObjectStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{
		client: client,
		config: config,
	}, nil
}

// Upload 上传文件
func (s *MinIOStorage) Upload(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.config.Bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download 下载文件
func (s *MinIOStorage) Download(ctx context.Context, objectName string) (io.Reader, error) {
	return s.client.GetObject(ctx, s.config.Bucket, objectName, minio.GetObjectOptions{})
}

// Delete 删除文件
func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.config.Bucket, objectName, minio.RemoveObjectOptions{})
}

// Exists 检查文件是否存在
func (s *MinIOStorage) Exists(ctx context.Context, objectName string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.config.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetURL 获取下载 URL（预签名）
func (s *MinIOStorage) GetURL(ctx context.Context, objectName string) (string, error) {
	// 生成 7 天有效的预签名 URL
	return s.client.PresignedGetObject(ctx, s.config.Bucket, objectName, 60*60*24*7, nil)
}
