package utils

import (
	"strings"
)

func GetSystemFileURL(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	baseUrl := GetBaseUrl()
	if baseUrl == "" {
		return path
	}

	baseUrl = strings.TrimSuffix(baseUrl, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseUrl + path
}

func GenerateFullURL(path string, storageType string) string {
	if storageType != "local" {
		return path
	}
	return GetSystemFileURL(path)
}

func GetFileFullURL(fileID string) string {
	baseUrl := GetBaseUrl()
	if baseUrl == "" {
		return "/f/" + fileID
	}
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	return baseUrl + "/f/" + fileID
}

func GetFileThumbnailFullURL(fileID string) string {
	baseUrl := GetBaseUrl()
	if baseUrl == "" {
		return "/t/" + fileID
	}
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	return baseUrl + "/t/" + fileID
}
