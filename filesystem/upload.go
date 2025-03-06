package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	fm_pb "github.com/nanoDFS/client-sdk/filesystem/proto/master"
	"github.com/nanoDFS/client-sdk/utils/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Upload returns fileId, userId, err
func (t *FileSystem) Upload(filePath string) (string, string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", "", err
	}
	file, _ := os.Open(filePath)
	masterAddr := config.LoadConfig().Master.Addr

	conn, err := grpc.NewClient(masterAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	fileId, userId := uuid.NewString(), uuid.NewString()

	client := fm_pb.NewFileMetadataServiceClient(conn)

	resp, err := client.UploadFile(context.Background(), &fm_pb.FileUploadReq{FileId: fileId, UserId: userId, Size: info.Size()})
	if err != nil || !resp.Success {
		return "", "", fmt.Errorf("failed to upload file: %v", err)
	}

	chunkSize := config.LoadConfig().Chunk.Size
	chunkCount := info.Size() / chunkSize

	wg := &sync.WaitGroup{}

	for i := range chunkCount {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == chunkCount-1 {
			end = info.Size()
		}

		go t.streamUpload(wg, file, start, end, resp.ChunkServers[i].Address, i, string(resp.GetAccessToken()))
	}

	wg.Wait()

	return fileId, userId, nil
}

func (t *FileSystem) streamUpload(wg *sync.WaitGroup, file *os.File, start int64, end int64, serverAddr string, chunkId int64, token string) {
	wg.Add(1)
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
	stream, err := client.Write(ctx)
	if err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}

	if _, err := file.Seek(start, io.SeekStart); err != nil {
		log.Fatalf("failed to seek file: %v", err)
	}

	buff := make([]byte, 1024)

	for {
		n, err := file.Read(buff)
		if err != nil && err.Error() != "EOF" {
			log.Fatalf("failed to read data: %v", err)
		}

		currentPos, _ := file.Seek(0, io.SeekCurrent)
		if n == 0 || currentPos > end {
			break
		}

		if err := stream.Send(&fs_pb.Payload{Data: buff[:n]}); err != nil {
			log.Fatalf("failed to send payload: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}

	log.Printf("Stream closed successfully: %v\n", resp)
}
