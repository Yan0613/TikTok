package ossClient

import (
	"github.com/Yan0613/TikTok/config"
	"github.com/Yan0613/TikTok/log/logger"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/credentials-go/credentials"
	"io"
	"mime/multipart"
)

type Credentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}

type CompletedPart struct {
	PartNumber int    // 分片编号
	ETag       string // 分片的ETag
}

type CredentialsProvider struct {
	cred credentials.Credential
}

func (credentials *Credentials) GetAccessKeyID() string {
	return credentials.AccessKeyId
}

func (credentials *Credentials) GetAccessKeySecret() string {
	return credentials.AccessKeySecret
}

func (credentials *Credentials) GetSecurityToken() string {
	return credentials.SecurityToken
}

func (defBuild CredentialsProvider) GetCredentials() oss.Credentials {
	id, _ := defBuild.cred.GetAccessKeyId()
	secret, _ := defBuild.cred.GetAccessKeySecret()
	token, _ := defBuild.cred.GetSecurityToken()

	return &Credentials{
		AccessKeyId:     *id,
		AccessKeySecret: *secret,
		SecurityToken:   *token,
	}
}

func NewStaticCredentialsProvider(credential credentials.Credential) CredentialsProvider {
	return CredentialsProvider{
		cred: credential,
	}
}

type MyBucket struct {
	*oss.Bucket
}

var myConfig = config.OSS()

func NewBucket(mode string) (*MyBucket, error) {
	config := new(credentials.Config).
		// 指定Credential类型，固定值为ecs_ram_role。
		SetType(myConfig.CredentialType).
		// （可选项）指定角色名称。如果不指定，OSS会自动获取角色。强烈建议指定角色名称，以降低请求次数。
		SetRoleName(myConfig.CredentialRoleName)

	ecsCredential, err := credentials.NewCredential(config)
	if err != nil {
		return nil, err
	}
	provider := NewStaticCredentialsProvider(ecsCredential)
	client, err := oss.New(myConfig.Endpoint[mode], "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		logger.Errorln("Failed to init OSS client:", err)
		return nil, err
	}
	// 存储空间名称
	bucket, err := client.Bucket(myConfig.BucketName)
	if err != nil {
		logger.Errorln("Failed to init bucket:", err)
		return nil, err
	}
	return &MyBucket{bucket}, nil
}

func (mb *MyBucket) UploadVideo(file *multipart.FileHeader, internalURL string) error {
	fileStream, err := file.Open()
	if err != nil {
		logger.Errorln("Failed to open multipart file:", internalURL, err)
		return err
	}
	defer fileStream.Close()

	imur, err := mb.InitiateMultipartUpload(internalURL)
	if err != nil {
		logger.Errorln("Failed to initiate multipart upload:", internalURL, err)
		return err
	}

	// 定义分片大小,此处为5MB
	partSize := int64(5 * 1024 * 1024)
	fileSize := file.Size

	// 计算分片数量，向上取整
	numParts := int((fileSize + partSize - 1) / partSize)

	var parts []oss.UploadPart

	// 分片上传
	for i := 0; i < numParts; i++ {
		start := int64(i) * partSize
		size := partSize
		if start+size > fileSize {
			size = fileSize - start
		}

		_, err = fileStream.Seek(start, io.SeekStart)
		if err != nil {
			logger.Errorln("Failed to seek file:", internalURL, err)
			mb.AbortMultipartUpload(imur)
			return err
		}

		part, err := mb.UploadPart(imur, fileStream, size, i+1)
		if err != nil {
			logger.Errorln("Error uploading part:", i+1, err)
			mb.AbortMultipartUpload(imur)
			return err
		}
		parts = append(parts, part)
	}

	// 完成分片上传
	cmur, err := mb.CompleteMultipartUpload(imur, parts)
	if err != nil {
		logger.Errorln("Failed to complete multipart upload:", internalURL, err)
		mb.AbortMultipartUpload(imur)
		return err
	}

	logger.Infof("File uploaded successfully: %s", cmur)
	return nil
}

const SignedURLExpiration = 600

func (mb *MyBucket) ObjectExternalURL(internalURL string) (signedURL string, err error) {
	// 生成用于下载的签名URL，并指定签名URL的有效时间为60秒。
	for i := 0; i < 5; i++ {
		signedURL, err = mb.SignURL(internalURL, oss.HTTPGet, SignedURLExpiration)
		if err != nil {
			logger.Errorln("Failed to get object URL:", internalURL, err)
			continue
		}
		return signedURL, nil
	}
	return "", err
}
