package logx

import (
	"context"

	"github.com/huaxr/framework/internal/define"
	"go.uber.org/zap"
)

// any-other fields need prepared should regulate defineFields size.
type defineFields [6]zap.Field

func (l *zLogger) getFields(ctx context.Context) defineFields {
	field := l.pool.Get().(defineFields)
	defer func() {
		l.pool.Put(field)
	}()

	field[0] = TraceId(ctx)
	field[1], field[2] = DurationPair(ctx)
	field[3] = CallFrom(ctx)
	field[4] = zap.String(define.HandlerPath.String(), "")
	field[5] = zap.Int64(define.HandlerExecutePeriod.String(), -1)
	return field
}
