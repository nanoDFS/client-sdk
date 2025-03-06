package download

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/charmbracelet/log"
	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (t *Downloader) stream(wg *sync.WaitGroup, serverAddr string, chunkId int64, token string) {
	defer wg.Done()

	conn, err := grpc.NewClient(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	md := metadata.Pairs(
		"auth", string(token),
		"chunk_id", fmt.Sprintf("%d", chunkId),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := fs_pb.NewFileStreamingServiceClient(conn)
	stream, err := client.Read(ctx, &fs_pb.ReadReq{})
	if err != nil {
		log.Errorf("failed to initiate read: %v", err)
		return
	}

	file, err := os.Create(path.Join(t.filePath, fmt.Sprintf("temp-%d", chunkId)))
	if err != nil {
		log.Errorf("failed to create temp file: %v", err)
		return
	}
	for {
		payload, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to read data: %v", err)
		}
		file.Write(payload.Data)
	}
	if err := stream.CloseSend(); err != nil {
		log.Errorf("failed to close: %v", err)
		return
	}
	log.Infof("Stream closed successfully for chunk id: %d\n", chunkId)
}
