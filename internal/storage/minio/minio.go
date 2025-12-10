package minio

import (
	"auth_service/internal/config"
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	cfg := config.App.Minio
	var err error

	MinioClient, err = minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = MinioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	exists, err := MinioClient.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.Fatalf("Failed to check bucket existence: %v", err)
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})

		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}

		policy := `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::` + cfg.Bucket + `/*"]
                }
            ]
        }`

		err = MinioClient.SetBucketPolicy(ctx, cfg.Bucket, policy)
		if err != nil {
			log.Printf("Warning: Failed to set bucket policy: %v", err)
		}
	}

	log.Println("Minio connected succesfully")

}
