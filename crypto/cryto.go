package crypto

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func BasicAuth(realm, pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:{SHA}%s", realm, sha)
}
