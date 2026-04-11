package pgxtrace_test

import (
	"context"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pgx-contrib/pgxtrace"
)

// contextKey is a private type used to verify context chaining.
type contextKey int

// mockQueryTracer implements only pgx.QueryTracer.
type mockQueryTracer struct {
	id         int
	startCalls int
	endCalls   int
}

func (m *mockQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData) context.Context {
	m.startCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockQueryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	m.endCalls++
}

// mockAllTracer implements all five pgx tracer interfaces.
type mockAllTracer struct {
	id                  int
	queryStartCalls     int
	queryEndCalls       int
	connectStartCalls   int
	connectEndCalls     int
	batchStartCalls     int
	batchQueryCalls     int
	batchEndCalls       int
	prepareStartCalls   int
	prepareEndCalls     int
	copyFromStartCalls  int
	copyFromEndCalls    int
}

func (m *mockAllTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData) context.Context {
	m.queryStartCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockAllTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
	m.queryEndCalls++
}

func (m *mockAllTracer) TraceConnectStart(ctx context.Context, _ pgx.TraceConnectStartData) context.Context {
	m.connectStartCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockAllTracer) TraceConnectEnd(_ context.Context, _ pgx.TraceConnectEndData) {
	m.connectEndCalls++
}

func (m *mockAllTracer) TraceBatchStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceBatchStartData) context.Context {
	m.batchStartCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockAllTracer) TraceBatchQuery(_ context.Context, _ *pgx.Conn, _ pgx.TraceBatchQueryData) {
	m.batchQueryCalls++
}

func (m *mockAllTracer) TraceBatchEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceBatchEndData) {
	m.batchEndCalls++
}

func (m *mockAllTracer) TracePrepareStart(ctx context.Context, _ *pgx.Conn, _ pgx.TracePrepareStartData) context.Context {
	m.prepareStartCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockAllTracer) TracePrepareEnd(_ context.Context, _ *pgx.Conn, _ pgx.TracePrepareEndData) {
	m.prepareEndCalls++
}

func (m *mockAllTracer) TraceCopyFromStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceCopyFromStartData) context.Context {
	m.copyFromStartCalls++
	return context.WithValue(ctx, contextKey(m.id), true)
}

func (m *mockAllTracer) TraceCopyFromEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceCopyFromEndData) {
	m.copyFromEndCalls++
}

var _ = Describe("CompositeQueryTracer", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	// -------------------------------------------------------------------------
	Describe("TraceQueryStart / TraceQueryEnd", func() {
		It("empty composite returns the same context without panic", func() {
			t := pgxtrace.CompositeQueryTracer{}
			result := t.TraceQueryStart(ctx, nil, pgx.TraceQueryStartData{})
			Expect(result).To(Equal(ctx))
			Expect(func() { t.TraceQueryEnd(ctx, nil, pgx.TraceQueryEndData{}) }).NotTo(Panic())
		})

		It("single tracer: TraceQueryStart and TraceQueryEnd each called once", func() {
			m := &mockQueryTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			t.TraceQueryStart(ctx, nil, pgx.TraceQueryStartData{})
			t.TraceQueryEnd(ctx, nil, pgx.TraceQueryEndData{})
			Expect(m.startCalls).To(Equal(1))
			Expect(m.endCalls).To(Equal(1))
		})

		It("multiple tracers: all called and contexts chained in order", func() {
			m1 := &mockQueryTracer{id: 1}
			m2 := &mockQueryTracer{id: 2}
			t := pgxtrace.CompositeQueryTracer{m1, m2}
			result := t.TraceQueryStart(ctx, nil, pgx.TraceQueryStartData{})
			Expect(m1.startCalls).To(Equal(1))
			Expect(m2.startCalls).To(Equal(1))
			// Both sentinel values must be present in the chained context.
			Expect(result.Value(contextKey(1))).To(BeTrue())
			Expect(result.Value(contextKey(2))).To(BeTrue())
		})
	})

	// -------------------------------------------------------------------------
	Describe("TraceConnectStart / TraceConnectEnd", func() {
		It("tracer not implementing ConnectTracer is skipped", func() {
			m := &mockQueryTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceConnectStart(ctx, pgx.TraceConnectStartData{})
			Expect(result).To(Equal(ctx))
			t.TraceConnectEnd(ctx, pgx.TraceConnectEndData{})
			// QueryTracer methods must not have been called.
			Expect(m.startCalls).To(Equal(0))
		})

		It("tracer implementing ConnectTracer is called", func() {
			m := &mockAllTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceConnectStart(ctx, pgx.TraceConnectStartData{})
			t.TraceConnectEnd(result, pgx.TraceConnectEndData{})
			Expect(m.connectStartCalls).To(Equal(1))
			Expect(m.connectEndCalls).To(Equal(1))
		})

		It("multiple full tracers: contexts chained", func() {
			m1 := &mockAllTracer{id: 1}
			m2 := &mockAllTracer{id: 2}
			t := pgxtrace.CompositeQueryTracer{m1, m2}
			result := t.TraceConnectStart(ctx, pgx.TraceConnectStartData{})
			Expect(result.Value(contextKey(1))).To(BeTrue())
			Expect(result.Value(contextKey(2))).To(BeTrue())
		})
	})

	// -------------------------------------------------------------------------
	Describe("TracePrepareStart / TracePrepareEnd", func() {
		It("tracer not implementing PrepareTracer is skipped", func() {
			m := &mockQueryTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TracePrepareStart(ctx, nil, pgx.TracePrepareStartData{})
			Expect(result).To(Equal(ctx))
			t.TracePrepareEnd(ctx, nil, pgx.TracePrepareEndData{})
			Expect(m.startCalls).To(Equal(0))
		})

		It("tracer implementing PrepareTracer is called", func() {
			m := &mockAllTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TracePrepareStart(ctx, nil, pgx.TracePrepareStartData{})
			t.TracePrepareEnd(result, nil, pgx.TracePrepareEndData{})
			Expect(m.prepareStartCalls).To(Equal(1))
			Expect(m.prepareEndCalls).To(Equal(1))
		})

		It("multiple full tracers: contexts chained", func() {
			m1 := &mockAllTracer{id: 1}
			m2 := &mockAllTracer{id: 2}
			t := pgxtrace.CompositeQueryTracer{m1, m2}
			result := t.TracePrepareStart(ctx, nil, pgx.TracePrepareStartData{})
			Expect(result.Value(contextKey(1))).To(BeTrue())
			Expect(result.Value(contextKey(2))).To(BeTrue())
		})
	})

	// -------------------------------------------------------------------------
	Describe("TraceBatchStart / TraceBatchQuery / TraceBatchEnd", func() {
		It("tracer not implementing BatchTracer is skipped", func() {
			m := &mockQueryTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceBatchStart(ctx, nil, pgx.TraceBatchStartData{})
			Expect(result).To(Equal(ctx))
			t.TraceBatchQuery(ctx, nil, pgx.TraceBatchQueryData{})
			t.TraceBatchEnd(ctx, nil, pgx.TraceBatchEndData{})
			Expect(m.startCalls).To(Equal(0))
		})

		It("tracer implementing BatchTracer is called", func() {
			m := &mockAllTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceBatchStart(ctx, nil, pgx.TraceBatchStartData{})
			t.TraceBatchQuery(result, nil, pgx.TraceBatchQueryData{})
			t.TraceBatchEnd(result, nil, pgx.TraceBatchEndData{})
			Expect(m.batchStartCalls).To(Equal(1))
			Expect(m.batchQueryCalls).To(Equal(1))
			Expect(m.batchEndCalls).To(Equal(1))
		})

		It("multiple full tracers: contexts chained", func() {
			m1 := &mockAllTracer{id: 1}
			m2 := &mockAllTracer{id: 2}
			t := pgxtrace.CompositeQueryTracer{m1, m2}
			result := t.TraceBatchStart(ctx, nil, pgx.TraceBatchStartData{})
			Expect(result.Value(contextKey(1))).To(BeTrue())
			Expect(result.Value(contextKey(2))).To(BeTrue())
		})
	})

	// -------------------------------------------------------------------------
	Describe("TraceCopyFromStart / TraceCopyFromEnd", func() {
		It("tracer not implementing CopyFromTracer is skipped", func() {
			m := &mockQueryTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceCopyFromStart(ctx, nil, pgx.TraceCopyFromStartData{})
			Expect(result).To(Equal(ctx))
			t.TraceCopyFromEnd(ctx, nil, pgx.TraceCopyFromEndData{})
			Expect(m.startCalls).To(Equal(0))
		})

		It("tracer implementing CopyFromTracer is called", func() {
			m := &mockAllTracer{id: 1}
			t := pgxtrace.CompositeQueryTracer{m}
			result := t.TraceCopyFromStart(ctx, nil, pgx.TraceCopyFromStartData{})
			t.TraceCopyFromEnd(result, nil, pgx.TraceCopyFromEndData{})
			Expect(m.copyFromStartCalls).To(Equal(1))
			Expect(m.copyFromEndCalls).To(Equal(1))
		})

		It("multiple full tracers: contexts chained", func() {
			m1 := &mockAllTracer{id: 1}
			m2 := &mockAllTracer{id: 2}
			t := pgxtrace.CompositeQueryTracer{m1, m2}
			result := t.TraceCopyFromStart(ctx, nil, pgx.TraceCopyFromStartData{})
			Expect(result.Value(contextKey(1))).To(BeTrue())
			Expect(result.Value(contextKey(2))).To(BeTrue())
		})
	})

	// -------------------------------------------------------------------------
	Describe("interface compliance", func() {
		It("satisfies pgx.QueryTracer", func() {
			var _ pgx.QueryTracer = pgxtrace.CompositeQueryTracer(nil)
		})

		It("satisfies pgx.BatchTracer", func() {
			var _ pgx.BatchTracer = pgxtrace.CompositeQueryTracer(nil)
		})

		It("satisfies pgx.ConnectTracer", func() {
			var _ pgx.ConnectTracer = pgxtrace.CompositeQueryTracer(nil)
		})

		It("satisfies pgx.PrepareTracer", func() {
			var _ pgx.PrepareTracer = pgxtrace.CompositeQueryTracer(nil)
		})

		It("satisfies pgx.CopyFromTracer", func() {
			var _ pgx.CopyFromTracer = pgxtrace.CompositeQueryTracer(nil)
		})
	})
})
