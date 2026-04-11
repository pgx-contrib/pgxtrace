package pgxtrace

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var (
	_ pgx.QueryTracer    = CompositeQueryTracer(nil)
	_ pgx.BatchTracer    = CompositeQueryTracer(nil)
	_ pgx.ConnectTracer  = CompositeQueryTracer(nil)
	_ pgx.PrepareTracer  = CompositeQueryTracer(nil)
	_ pgx.CopyFromTracer = CompositeQueryTracer(nil)
)

// CompositeQueryTracer chains multiple pgx tracers into one. Every element
// must implement [pgx.QueryTracer]; elements that also implement
// [pgx.ConnectTracer], [pgx.BatchTracer], [pgx.PrepareTracer], or
// [pgx.CopyFromTracer] are automatically dispatched for those operations.
// Start methods are called in registration order; End methods are called in
// reverse order, mirroring the stack-unwinding semantics of defer.
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
	for i := len(t) - 1; i >= 0; i-- {
		if tracer, ok := t[i].(pgx.ConnectTracer); ok {
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
	for i := len(t) - 1; i >= 0; i-- {
		if tracer, ok := t[i].(pgx.PrepareTracer); ok {
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
	for i := len(t) - 1; i >= 0; i-- {
		t[i].TraceQueryEnd(ctx, conn, data)
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
	for i := len(t) - 1; i >= 0; i-- {
		if tracer, ok := t[i].(pgx.BatchTracer); ok {
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
	for i := len(t) - 1; i >= 0; i-- {
		if tracer, ok := t[i].(pgx.CopyFromTracer); ok {
			tracer.TraceCopyFromEnd(ctx, conn, data)
		}
	}
}
