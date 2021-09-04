package utils

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// CheckPasswordHash checks a password hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// RandString generates a random string. Taken from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go (RandStringBytesMaskImprSrcSB())
func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	loc, _ := time.LoadLocation("UTC")
	t = t.In(loc)

	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func GetLog() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.Development = true
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	cfg.Encoding = "console" // "console" or "json"
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.Level.SetLevel(zap.DebugLevel)

	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("Could not build logger, err: %v", err)
	}
	defer logger.Sync() // Flushes buffer, if any
	log := logger.Sugar()

	return log
}
