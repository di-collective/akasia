package utils

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

func PutFile(BucketName, Endpoint, KeyID, SecretKey, tempFileName, newFilename string) error {
	client, err := oss.New(Endpoint, KeyID, SecretKey)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(newFilename, tempFileName)
	if err != nil {
		return err
	}
	return err
}
