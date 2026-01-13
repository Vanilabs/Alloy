package socket

import (
	"context"

	"time"

	"github.com/redis/go-redis/v9"
)

type SocketTracker struct {
	RDB *redis.Client
}

func NewSocketTracker(rdb *redis.Client) *SocketTracker {
	return &SocketTracker{
		RDB: rdb,
	}
}

func (t *SocketTracker) AddSocketForUser(userID, socketID string) error {
	ctx := context.Background()

	pipe := t.RDB.TxPipeline()

	pipe.SAdd(ctx, "user_sockets:"+userID, socketID)
	pipe.HSet(ctx, "socket:"+socketID, "userID", userID)

	// 👇 THIS IS THE KEY FIX
	pipe.Expire(ctx, "socket:"+socketID, 30*time.Second)

	_, err := pipe.Exec(ctx)
	return err
}


func (t *SocketTracker) RemoveSocket(socketID string) error {
	ctx := context.Background()

	userID, err := t.RDB.HGet(ctx, "socket:"+socketID, "userID").Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	t.RDB.SRem(ctx, "user_sockets:"+userID, socketID)
	t.RDB.Del(ctx, "socket:"+socketID)

	return nil
}


func (t *SocketTracker) GetUserSocketsExcept(userID, excludeSocketID string) ([]string, error) {

	sockets, err := t.GetUserSockets(userID)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(sockets))
	for _, sid := range sockets {
		if sid != excludeSocketID {
			result = append(result, sid)
		}
	}

	return result, nil
}

func (t *SocketTracker) GetUserSockets(userID string) ([]string, error) {
	ctx := context.Background()

	socketIDs, err := t.RDB.SMembers(ctx, "user_sockets:"+userID).Result()
	if err != nil {
		return nil, err
	}

	if len(socketIDs) == 0 {
		return nil, nil
	}

	valid := make([]string, 0, len(socketIDs))

	for _, sid := range socketIDs {
		exists, err := t.RDB.Exists(ctx, "socket:"+sid).Result()
		if err != nil {
			continue
		}

		if exists == 1 {
			valid = append(valid, sid)
		} else {
			t.RDB.SRem(ctx, "user_sockets:"+userID, sid)
		}
	}

	return valid, nil
}


func (t *SocketTracker) GetUserBySocket(socketID string) (string, error) {
	return t.RDB.HGet(context.Background(), "socket_owner:"+socketID, "userID").Result()
}
