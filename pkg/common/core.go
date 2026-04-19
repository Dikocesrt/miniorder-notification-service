package common

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	HEADER_TRACE_ID = "traceID"
	LOCAL_TRACE_ID  = "TRACE_ID"
)

func GetTraceID(ctx *fiber.Ctx) string {
	// From Header
	traceId := ctx.Get(HEADER_TRACE_ID)
	if traceId != "" {
		return traceId
	}

	// From Locals
	traceId, _ = ctx.Locals(LOCAL_TRACE_ID).(string)
	if traceId != "" {
		return traceId
	}

	// Generate New Trace ID
	traceId = CreateTraceID()
	ctx.Locals(LOCAL_TRACE_ID, traceId)
	return traceId
}

func CreateRandomHex(n int, charType string) string {
	uid := uuid.New()

	charset := ""
	switch charType {
	case "N":
		charset = "0123456789"
	case "A":
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "AN":
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	default:
		return strings.ToUpper(uid.String()[:n])
	}
	result := make([]byte, n)
	for i := range result {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result)
}

func CreateTraceID() string {
	randHex := CreateRandomHex(5, "AN")
	nowFormat := time.Now().Format("060102150405.000")
	return strings.Replace(nowFormat, ".", "", 1) + randHex
}

func ExecJob(worker ...func(wg *sync.WaitGroup)) {
	var wg sync.WaitGroup
	for _, work := range worker {
		wg.Add(1)
		go work(&wg)
	}
	wg.Wait()
}

func CompactJSON(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return string(data)
	}
	return buf.String()
}
