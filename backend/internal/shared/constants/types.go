package constants

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/gocql/gocql"

)

type DataStores struct {
	Redis   *redis.Client
	PostGres       *gorm.DB
	Cassandra 	*gocql.Session
}

type MessageChannelType string
type MessageChannelRole string


const (
	ChannelDM    MessageChannelType = "dm"
	ChannelGroup MessageChannelType = "group"

	GroupAdmin MessageChannelRole = "admin"
	GroupMember MessageChannelRole = "member"
)