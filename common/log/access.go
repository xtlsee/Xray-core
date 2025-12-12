package log

import (
	"context"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/xtls/xray-core/common/serial"
)

type logKey int

const (
	accessMessageKey logKey = iota
)

type AccessStatus string

const (
	AccessAccepted = AccessStatus("accepted")
	AccessRejected = AccessStatus("rejected")
)

type AccessMessage struct {
	From        interface{}
	To          interface{}
	Status      AccessStatus
	Reason      interface{}
	Email       string
	Detour      string
	UpCounter   *atomic.Int64
	DownCounter *atomic.Int64
}

func (m *AccessMessage) String() string {
	builder := strings.Builder{}
	builder.WriteString("from")
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.From))
	builder.WriteByte(' ')
	builder.WriteString(string(m.Status))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.To))

	if len(m.Detour) > 0 {
		builder.WriteString(" [")
		builder.WriteString(m.Detour)
		builder.WriteByte(']')
	}

	if reason := serial.ToString(m.Reason); len(reason) > 0 {
		builder.WriteString(" ")
		builder.WriteString(reason)
	}

	if len(m.Email) > 0 {
		builder.WriteString(" email: ")
		builder.WriteString(m.Email)
	}

	if m.UpCounter != nil {
		builder.WriteString(" upload: ")
		builder.WriteString(strconv.FormatInt(m.UpCounter.Load(), 10))
	}
	if m.DownCounter != nil {
		builder.WriteString(" download: ")
		builder.WriteString(strconv.FormatInt(m.DownCounter.Load(), 10))
	}
	return builder.String()
}

func ContextWithAccessMessage(ctx context.Context, accessMessage *AccessMessage) context.Context {
	if accessMessage != nil {
		if accessMessage.UpCounter == nil {
			accessMessage.UpCounter = new(atomic.Int64)
		}
		if accessMessage.DownCounter == nil {
			accessMessage.DownCounter = new(atomic.Int64)
		}
	}

	return context.WithValue(ctx, accessMessageKey, accessMessage)
}

func AccessMessageFromContext(ctx context.Context) *AccessMessage {
	if accessMessage, ok := ctx.Value(accessMessageKey).(*AccessMessage); ok {
		return accessMessage
	}
	return nil
}

func RecordFromContext(ctx context.Context) {
	if accessMessage := AccessMessageFromContext(ctx); accessMessage != nil {
		Record(accessMessage)
	}
}
