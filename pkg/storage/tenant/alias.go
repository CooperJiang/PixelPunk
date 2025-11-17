package tenant

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"strings"
	"sync"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
)

// Simple in-process cache for userID -> alias
var aliasCache sync.Map // key: uint (userID), value: string

// ResolveAlias returns a stable per-user alias, creating one lazily if missing.
func ResolveAlias(userID uint) (string, error) {
	if userID == 0 {
		if v, ok := aliasCache.Load(userID); ok {
			if s, ok2 := v.(string); ok2 && s != "" {
				return s, nil
			}
		}
		aliasCache.Store(userID, "system")
		return "system", nil
	}
	if v, ok := aliasCache.Load(userID); ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			return s, nil
		}
	}
	db := database.GetDB()
	if db == nil {
		return "", sql.ErrConnDone
	}
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	alias := strings.TrimSpace(user.PathAlias)
	if alias == "" {
		// generate and persist
		for i := 0; i < 5; i++ {
			cand := genAlias(10)
			// ensure uniqueness
			var cnt int64
			if err := db.Model(&models.User{}).Where("path_alias = ?", cand).Count(&cnt).Error; err != nil {
				continue
			}
			if cnt == 0 {
				if err := db.Model(&models.User{}).Where("id = ?", userID).Update("path_alias", cand).Error; err == nil {
					alias = cand
					break
				}
			}
		}
		if alias == "" {
			// fallback to derived but non-empty alias to avoid panic paths
			alias = genAlias(10)
		}
	}
	aliasCache.Store(userID, alias)
	return alias, nil
}

// genAlias returns a base32 (lower) string with given length, using crypto/rand bytes.
func genAlias(n int) string {
	// Generate 16 random bytes and base32 encode, then trim to n
	var b [16]byte
	_, _ = rand.Read(b[:])
	// use no padding and lowercase
	enc := base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)
	s := strings.ToLower(enc.EncodeToString(b[:]))
	if len(s) < n {
		return s
	}
	return s[:n]
}

// shardPrefix builds a two-level shard directory from alias (first two chars).
func shardPrefix(alias string) (string, string) {
	a1, a2 := "x", "x"
	if len(alias) >= 1 {
		a1 = string(alias[0])
	}
	if len(alias) >= 2 {
		a2 = string(alias[1])
	}
	return a1, a2
}

// BuildObjectKey builds cloud object key for original files using alias sharding.
// files/{a1}{a2}/u_{alias}/(folder)/file
func BuildObjectKey(userID uint, folderPath, fileName string) (string, error) {
	// 游客：固定放在 files/guest 下，避免别名与分片
	if userID == 0 {
		return "files/guest/" + fileName, nil
	}
	alias, err := ResolveAlias(userID)
	if err != nil {
		return "", err
	}
	a1, a2 := shardPrefix(alias)
	base := "files/" + a1 + a2 + "/" + alias
	folderPath = strings.TrimPrefix(folderPath, "/")
	if folderPath != "" {
		return base + "/" + folderPath + "/" + fileName, nil
	}
	return base + "/" + fileName, nil
}

// BuildThumbObjectKey builds cloud object key for thumbnails using alias sharding.
// thumbnails/{a1}{a2}/u_{alias}/(folder)/file
func BuildThumbObjectKey(userID uint, folderPath, fileName string) (string, error) {
	// 游客缩略图：固定放在 thumbnails/guest 下
	if userID == 0 {
		return "thumbnails/guest/" + fileName, nil
	}
	alias, err := ResolveAlias(userID)
	if err != nil {
		return "", err
	}
	a1, a2 := shardPrefix(alias)
	base := "thumbnails/" + a1 + a2 + "/" + alias
	folderPath = strings.TrimPrefix(folderPath, "/")
	if folderPath != "" {
		return base + "/" + folderPath + "/" + fileName, nil
	}
	return base + "/" + fileName, nil
}
