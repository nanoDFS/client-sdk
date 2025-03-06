package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (t *Uploader) stream(wg *sync.WaitGroup, file *os.File, start int64, end int64, serverAddr string, chunkId int64, token string) {
	defer wg.Done()
	conn, err := grpc.NewClient(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		log.Errorf("%v", err)
		return
	}
	defer conn.Close()
	md := metadata.Pairs(
		"auth", string(token),
		"chunk_id", fmt.Sprintf("%d", chunkId),
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := fs_pb.NewFileStreamingServiceClient(conn)
	stream, err := client.Write(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		log.Errorf("%v", err)
		return
	}
	buff := make([]byte, 1024)
	for {
		n, err := file.Read(buff)
		if err != nil && err.Error() != "EOF" {
			log.Errorf("%v", err)
			return
		}
		currentPos, _ := file.Seek(0, io.SeekCurrent)
		if n == 0 || currentPos > end {
			break
		}
		if err := stream.Send(&fs_pb.Payload{Data: buff[:n]}); err != nil {
			log.Errorf("%v", err)
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("Stream closed successfully: %v\n", resp)
}
