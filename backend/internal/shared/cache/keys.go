package cache

import "fmt"

func GetMagicLinkKey(token string) string {
	return fmt.Sprintf("magic_link:%s", token)
}
