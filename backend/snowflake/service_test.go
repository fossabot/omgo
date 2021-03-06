package main

import (
	"context"
	"testing"

	"github.com/master-g/omgo/kit/utils"
	pb "github.com/master-g/omgo/proto/grpc/snowflake"
	"google.golang.org/grpc"
)

const (
	testKey = "test_key"
)

var address string

func init() {
	address = utils.GetLocalIP() + ":40001"
}

func TestCasDelay(t *testing.T) {
	casDelay()
}

func TestSnowflake(t *testing.T) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Next(context.Background(), &pb.Snowflake_Key{Name: testKey})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Log(r.Value)
}

func BenchmarkSnowflake(b *testing.B) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	for i := 0; i < b.N; i++ {
		// Contact the server and print out its response.
		_, err := c.Next(context.Background(), &pb.Snowflake_Key{Name: testKey})
		if err != nil {
			b.Fatalf("Could not get next value: %v", err)
		}
	}
}

func TestSnowflakeNext2(t *testing.T) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.Next2(context.Background(), &pb.Snowflake_Param{Name: testKey, Step: 100})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Log(r.Value)

	r, err = c.Next2(context.Background(), &pb.Snowflake_Param{Name: testKey})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Log(r.Value)
}

func BenchmarkSnowflakeNext2(b *testing.B) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	for i := 0; i < b.N; i++ {
		// Contact the server and print out its response.
		_, err := c.Next2(context.Background(), &pb.Snowflake_Param{Name: testKey})
		if err != nil {
			b.Fatalf("Could not get next value: %v", err)
		}
	}
}

func TestSnowflakeUUID(t *testing.T) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	// Contact the server and print out its response.
	r, err := c.GetUUID(context.Background(), &pb.Snowflake_NullRequest{})
	if err != nil {
		t.Fatalf("could not get next value: %v", err)
	}
	t.Logf("%b", r.Uuid)
}

func BenchmarkSnowflakeUUID(b *testing.B) {
	// Setup a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	c := pb.NewSnowflakeServiceClient(conn)

	for i := 0; i < b.N; i++ {
		// Contact the server and print out its response.
		_, err := c.GetUUID(context.Background(), &pb.Snowflake_NullRequest{})
		if err != nil {
			b.Fatalf("could not get uuid: %v", err)
		}
	}
}
