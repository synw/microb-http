package csrf

import (
	"crypto/rand"
	"fmt"
	"github.com/synw/microb/redis"
	"github.com/synw/terr"
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func VerifyToken(token string) *terr.Trace {
	_, err := redis.GetKey(token)
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	return nil
}

func GetToken() (string, *terr.Trace) {
	conn := redis.GetConn()
	defer conn.Close()
	token := randToken()
	_, err := conn.Do("SET", token, "")
	if err != nil {
		tr := terr.New(err)
		return token, tr
	}
	ttl := 900 // the tokens will live for 15 minutes
	_, err = conn.Do("EXPIRE", token, ttl)
	if err != nil {
		tr := terr.New(err)
		return token, tr
	}
	return token, nil
}
