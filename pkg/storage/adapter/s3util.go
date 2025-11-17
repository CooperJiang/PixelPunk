package adapter

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// s3ExtractHost trims scheme and returns host for an endpoint like https://host or host.
func s3ExtractHost(endpoint string) string {
	h := strings.TrimSpace(endpoint)
	h = strings.TrimPrefix(h, "https://")
	h = strings.TrimPrefix(h, "http://")
	return h
}

// s3DefaultHostForRegion returns AWS S3 host for region.
func s3DefaultHostForRegion(region string) string {
	if strings.TrimSpace(region) == "" || strings.EqualFold(region, "us-east-1") {
		return "s3.amazonaws.com"
	}
	return "s3." + region + ".amazonaws.com"
}

// s3BuildURL builds public URL for S3-style storage.
// If customDomain provided, always use scheme://customDomain/path
// Otherwise choose between path-style and virtual-host-style.
func s3BuildURL(scheme, bucket, endpointHost, path string, usePathStyle bool, customDomain string) string {
	if cd := strings.TrimSpace(customDomain); cd != "" {
		return scheme + "://" + strings.TrimSuffix(cd, "/") + "/" + encodePathSegments(path)
	}
	if usePathStyle {
		return scheme + "://" + endpointHost + "/" + bucket + "/" + encodePathSegments(path)
	}
	return scheme + "://" + bucket + "." + endpointHost + "/" + encodePathSegments(path)
}

// s3MapACL maps simple string to S3 canned ACL.
func s3MapACL(ac string) (types.ObjectCannedACL, bool) {
	switch strings.ToLower(strings.TrimSpace(ac)) {
	case "public-read":
		return types.ObjectCannedACLPublicRead, true
	case "private":
		return types.ObjectCannedACLPrivate, true
	default:
		return "", false
	}
}
