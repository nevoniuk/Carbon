/*
package clickhouse wraps the native clickhouse driver and adds retries to
alleviate network errors such as:
write tcp 10.20.1.69:49014->54.184.104.213:9440: write: broken pipe
*/

package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type (
	// Conn is a ClickHouse driver connection object that implement retries.
	Conn interface {
		Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
		QueryRow(ctx context.Context, query string, args ...interface{}) Row
		PrepareBatch(ctx context.Context, query string) (Batch, error)
		Exec(ctx context.Context, query string, args ...interface{}) error
		Ping(context.Context) error
		Close() error
	}

	Row interface {
		Err() error
		Scan(dest ...interface{}) error
		ScanStruct(dest interface{}) error
	}

	Rows interface {
		Next() bool
		Scan(dest ...interface{}) error
		ScanStruct(dest interface{}) error
		ColumnTypes() []driver.ColumnType
		Totals(dest ...interface{}) error
		Columns() []string
		Close() error
		Err() error
	}

	Batch interface {
		Abort() error
		Append(v ...interface{}) error
		AppendStruct(v interface{}) error
		Column(int) driver.BatchColumn
		Send() error
	}

	BatchColumn interface {
		Append(interface{}) error
	}

	// client implements the Client interface.
	client struct {
		driver.Conn
	}

	// rows wraps the native clickhouse driver's Rows.
	rows struct {
		driver.Rows
	}

	// row wraps the native clickhouse driver's Row.
	row struct {
		driver.Row
	}

	// batch wraps the native clickhouse driver's Batch.
	batch struct {
		driver.Batch
	}
)

// New returns a new Client.
func New(conn driver.Conn) Conn {
	return &client{conn}
}

// Query executes a query and returns Rows.
func (c *client) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var rs driver.Rows
	var err error
	err = WithRetries(func() error {
		rs, err = c.Conn.Query(ctx, query, args...)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &rows{rs}, nil
}

// QueryRow executes a query and returns a Row.
func (c *client) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	var r driver.Row
	WithRetries(func() error {
		r = c.Conn.QueryRow(ctx, query, args...)
		return r.Err()
	})
	return &row{r}
}

// PrepareBatch prepares a batch query.
func (c *client) PrepareBatch(ctx context.Context, query string) (Batch, error) {
	var b driver.Batch
	var err error
	err = WithRetries(func() error {
		b, err = c.Conn.PrepareBatch(ctx, query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &batch{b}, nil
}

// Exec executes a query.
func (c *client) Exec(ctx context.Context, query string, args ...interface{}) error {
	return WithRetries(func() error {
		return c.Conn.Exec(ctx, query, args...)
	})
}

//Ping pings the server.
func (c *client) Ping(ctx context.Context) error {
	return WithRetries(func() error {
		return c.Conn.Ping(ctx)
	})
}

// Scan scans the row into dest.
func (r *row) Scan(dest ...interface{}) error {
	return WithRetries(func() error {
		return r.Row.Scan(dest...)
	})
}

// ScanStruct calls driver.Rows.ScanStruct with retries.
func (r *row) ScanStruct(dest interface{}) error {
	return WithRetries(func() error {
		return r.Row.ScanStruct(dest)
	})
}

// Scan scans the row into dest.
func (r *rows) Scan(dest ...interface{}) error {
	return WithRetries(func() error {
		return r.Rows.Scan(dest...)
	})
}

// ScanStruct scans the row into dest.
func (r *rows) ScanStruct(dest interface{}) error {
	return WithRetries(func() error {
		return r.Rows.ScanStruct(dest)
	})
}

// Totals runs Totals with retries.
func (r *rows) Totals(dest ...interface{}) error {
	return WithRetries(func() error {
		return r.Rows.Totals(dest...)
	})
}

// Send sends the batch.
func (b *batch) Send() error {
	return WithRetries(func() error {
		return b.Batch.Send()
	})
}