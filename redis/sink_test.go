package redis

import (
	"errors"
	"reflect"
	"testing"

	"github.com/syllabix/logger/redis/internal/mocks"
)

func TestSync_Write(t *testing.T) {

	msg1 := "INFO source_host=my-k8-pod message=hello world, we miss you"
	msg2 := "ERROR source_host=my-k8-pod message=something seriously is not good"

	conn := new(mocks.Conn)
	conn.On("Do", "RPUSH", []interface{}{"logstash.stg.bolcom", msg1}).Return("ok", nil)
	conn.On("Do", "RPUSH", []interface{}{"logstash.stg.bolcom", msg2}).Return(nil, errors.New("something is not right"))
	conn.On("Close").Return(nil)

	mpool := new(mocks.Pool)
	mpool.On("Get").Return(conn)

	type fields struct {
		pool Pool
		key  string
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantN    int
		wantErr  bool
		asserter func(t *testing.T)
	}{
		{
			name: "ok write",
			fields: fields{
				pool: mpool,
				key:  "logstash.stg.bolcom",
			},
			args: args{
				p: []byte(msg1),
			},
			wantN: 59,
			asserter: func(t *testing.T) {
				conn.AssertCalled(t, "Do", "RPUSH", []interface{}{"logstash.stg.bolcom", msg1})
				conn.AssertCalled(t, "Close")
			},
		},
		{
			name: "ok write",
			fields: fields{
				pool: mpool,
				key:  "logstash.stg.bolcom",
			},
			args: args{
				p: []byte(msg2),
			},
			wantN:   0,
			wantErr: true,
			asserter: func(t *testing.T) {
				conn.AssertCalled(t, "Do", "RPUSH", []interface{}{"logstash.stg.bolcom", msg2})
				conn.AssertCalled(t, "Close")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sink{
				pool: tt.fields.pool,
				key:  tt.fields.key,
			}
			gotN, err := s.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sync.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Sync.Write() = %v, want %v", gotN, tt.wantN)
			}
			tt.asserter(t)
		})
	}
}

func TestSync_Close(t *testing.T) {

	okpool := new(mocks.Pool)
	okpool.On("Close").Return(nil)

	errPool := new(mocks.Pool)
	errPool.On("Close").Return(errors.New("yeah - that sucked"))

	type fields struct {
		pool Pool
		key  string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		asserter func(t *testing.T)
	}{
		{
			name: "close ok",
			fields: fields{
				pool: okpool,
				key:  "logstash.stg.bolcom",
			},
			asserter: func(t *testing.T) {
				okpool.AssertNumberOfCalls(t, "Close", 1)
			},
		},
		{
			name: "close error",
			fields: fields{
				pool: errPool,
				key:  "logstash.stg.bolcom",
			},
			asserter: func(t *testing.T) {
				errPool.AssertNumberOfCalls(t, "Close", 1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sink{
				pool: tt.fields.pool,
				key:  tt.fields.key,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Sync.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSync(t *testing.T) {

	okpool := new(mocks.Pool)

	type args struct {
		key  string
		pool Pool
	}
	tests := []struct {
		name string
		args args
		want *Sink
	}{
		{
			name: "key and pool",
			args: args{
				key:  "logstash.pro.bolcom",
				pool: okpool,
			},
			want: &Sink{
				pool: okpool,
				key:  "logstash.pro.bolcom",
			},
		},
		{
			name: "nil pool",
			args: args{
				key:  "logstash.pro.bolcom",
				pool: nil,
			},
			want: &Sink{
				pool: nil,
				key:  "logstash.pro.bolcom",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSink(tt.args.key, tt.args.pool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSync() = %v, want %v", got, tt.want)
			}
		})
	}
}
