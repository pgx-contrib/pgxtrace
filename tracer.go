package pgxtrace

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var (
	_ pgx.QueryTracer    = (*CompositeQueryTracer)(nil)
	_ pgx.BatchTracer    = (*CompositeQueryTracer)(nil)
	_ pgx.ConnectTracer  = (*CompositeQueryTracer)(nil)
	_ pgx.PrepareTracer  = (*CompositeQueryTracer)(nil)
	_ pgx.CopyFromTracer = (*CompositeQueryTracer)(nil)
)

// CompositeQueryTracer represent a composite query tracer
type CompositeQueryTracer []pgx.QueryTracer

// TraceConnectStart implements pgx.ConnectTracer.
func (t CompositeQueryTracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	for _, item := range t {
		if tracer, ok := item.(pgx.ConnectTracer); ok {
			ctx = tracer.TraceConnectStart(ctx, data)
		}
	}

	return ctx
}

// TraceConnectEnd implements pgx.ConnectTracer.
func (t CompositeQueryTracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	for _, item := range t {
		if tracer, ok := item.(pgx.ConnectTracer); ok {
			tracer.TraceConnectEnd(ctx, data)
		}
	}
}

// TracePrepareStart implements pgx.PrepareTracer.
func (t CompositeQueryTracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	for _, item := range t {
		if tracer, ok := item.(pgx.PrepareTracer); ok {
			ctx = tracer.TracePrepareStart(ctx, conn, data)
		}
	}

	return ctx
}

// TracePrepareEnd implements pgx.PrepareTracer.
func (t CompositeQueryTracer) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	for _, item := range t {
		if tracer, ok := item.(pgx.PrepareTracer); ok {
			tracer.TracePrepareEnd(ctx, conn, data)
		}
	}
}

// TraceQueryStart implements pgx.QueryTracer.
func (t CompositeQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range t {
		ctx = tracer.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

// TraceQueryEnd implements pgx.QueryTracer.
func (t CompositeQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range t {
		tracer.TraceQueryEnd(ctx, conn, data)
	}
}

// TraceBatchStart implements pgx.BatchTracer.
func (t CompositeQueryTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	for _, item := range t {
		if tracer, ok := item.(pgx.BatchTracer); ok {
			ctx = tracer.TraceBatchStart(ctx, conn, data)
		}
	}

	return ctx
}

// TraceBatchQuery implements pgx.BatchTracer.
func (t CompositeQueryTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	for _, item := range t {
		if tracer, ok := item.(pgx.BatchTracer); ok {
			tracer.TraceBatchQuery(ctx, conn, data)
		}
	}
}

// TraceBatchEnd implements pgx.BatchTracer.
func (t CompositeQueryTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	for _, item := range t {
		if tracer, ok := item.(pgx.BatchTracer); ok {
			tracer.TraceBatchEnd(ctx, conn, data)
		}
	}
}

// TraceCopyFromStart implements pgx.CopyFromTracer.
func (t CompositeQueryTracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	for _, item := range t {
		if tracer, ok := item.(pgx.CopyFromTracer); ok {
			ctx = tracer.TraceCopyFromStart(ctx, conn, data)
		}
	}

	return ctx
}

// TraceCopyFromEnd implements pgx.CopyFromTracer.
func (t CompositeQueryTracer) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
	for _, item := range t {
		if tracer, ok := item.(pgx.CopyFromTracer); ok {
			tracer.TraceCopyFromEnd(ctx, conn, data)
		}
	}
}
