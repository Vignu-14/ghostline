package config

import "strings"

type StorageConfig struct {
	SupabaseURL        string
	SupabaseServiceKey string
	BucketName         string
}

func (c StorageConfig) Enabled() bool {
	return strings.TrimSpace(c.SupabaseURL) != "" &&
		strings.TrimSpace(c.SupabaseServiceKey) != "" &&
		strings.TrimSpace(c.BucketName) != ""
}
