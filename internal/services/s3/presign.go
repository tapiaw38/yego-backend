package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client holds S3 credentials and bucket config.
type Client struct {
	region    string
	bucket    string
	accessKey string
	secretKey string
}

// NewClient creates a new S3 presign client.
func NewClient(region, bucket, accessKey, secretKey string) *Client {
	return &Client{
		region:    region,
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// IsConfigured returns true if all required credentials are set.
func (c *Client) IsConfigured() bool {
	return c.region != "" && c.bucket != "" && c.accessKey != "" && c.secretKey != ""
}

// PresignPut generates an AWS S3 presigned PUT URL (SigV4) valid for the given duration.
// Returns (upload_url, public_url, error).
func (c *Client) PresignPut(key string, expires time.Duration) (string, string, error) {
	now := time.Now().UTC()
	datetime := now.Format("20060102T150405Z")
	date := now.Format("20060102")

	host := fmt.Sprintf("%s.s3.%s.amazonaws.com", c.bucket, c.region)
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", date, c.region)
	credential := fmt.Sprintf("%s/%s", c.accessKey, scope)
	expiresStr := fmt.Sprintf("%d", int(expires.Seconds()))

	canonicalQueryString := strings.Join([]string{
		"X-Amz-Algorithm=AWS4-HMAC-SHA256",
		"X-Amz-Credential=" + url.QueryEscape(credential),
		"X-Amz-Date=" + datetime,
		"X-Amz-Expires=" + expiresStr,
		"X-Amz-SignedHeaders=host",
	}, "&")

	encodedKey := pathEscape(key)

	canonicalRequest := strings.Join([]string{
		"PUT",
		"/" + encodedKey,
		canonicalQueryString,
		"host:" + host + "\n",
		"host",
		"UNSIGNED-PAYLOAD",
	}, "\n")

	hashedCanonical := sha256Hex(canonicalRequest)

	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		datetime,
		scope,
		hashedCanonical,
	}, "\n")

	signingKey := hmacSHA256(
		hmacSHA256(
			hmacSHA256(
				hmacSHA256([]byte("AWS4"+c.secretKey), date),
				c.region,
			),
			"s3",
		),
		"aws4_request",
	)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	uploadURL := fmt.Sprintf("https://%s/%s?%s&X-Amz-Signature=%s",
		host, encodedKey, canonicalQueryString, signature)
	publicURL := fmt.Sprintf("https://%s/%s", host, encodedKey)

	return uploadURL, publicURL, nil
}

// DeleteObject deletes an object from S3 by key.
func (c *Client) DeleteObject(key string) error {
	now := time.Now().UTC()
	datetime := now.Format("20060102T150405Z")
	date := now.Format("20060102")

	host := fmt.Sprintf("%s.s3.%s.amazonaws.com", c.bucket, c.region)
	encodedKey := "/" + pathEscape(key)
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", date, c.region)
	payloadHash := sha256Hex("")

	canonicalHeaders := "host:" + host + "\nx-amz-date:" + datetime + "\n"
	signedHeaders := "host;x-amz-date"
	canonicalRequest := strings.Join([]string{"DELETE", encodedKey, "", canonicalHeaders, signedHeaders, payloadHash}, "\n")
	stringToSign := strings.Join([]string{"AWS4-HMAC-SHA256", datetime, scope, sha256Hex(canonicalRequest)}, "\n")

	signingKey := hmacSHA256(hmacSHA256(hmacSHA256(hmacSHA256([]byte("AWS4"+c.secretKey), date), c.region), "s3"), "aws4_request")
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
	auth := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", c.accessKey, scope, signedHeaders, signature)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s%s", host, encodedKey), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Host", host)
	req.Header.Set("x-amz-date", datetime)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("s3 delete failed: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// pathEscape encodes a path preserving forward slashes.
func pathEscape(s string) string {
	parts := strings.Split(s, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}
