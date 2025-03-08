package download

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/nanoDFS/client-sdk/crypto"

	"github.com/charmbracelet/log"
	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type streamOpts struct {
	wg         *sync.WaitGroup
	serverAddr string
	chunkId    int64
	token      string
	key        crypto.CryptoKey
}

func (t *Downloader) stream(opts streamOpts) {
	defer opts.wg.Done()

	conn, err := grpc.NewClient(opts.serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	md := metadata.Pairs(
		"auth", string(opts.token),
		"chunk_id", fmt.Sprintf("%d", opts.chunkId),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := fs_pb.NewFileStreamingServiceClient(conn)
	stream, err := client.Read(ctx, &fs_pb.ReadReq{})
	if err != nil {
		log.Errorf("failed to initiate read: %v", err)
		return
	}

	file, err := os.Create(path.Join(t.filePath, fmt.Sprintf("temp-%d", opts.chunkId)))
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
		decryptedData, err := crypto.NewDecryptor().Decrypt(payload.Data, opts.key.Nonce, opts.key.Key)
		if err != nil {
			log.Errorf("failed to decrypt: %v", err)
			break
		}
		file.Write(decryptedData)
	}
	if err := stream.CloseSend(); err != nil {
		log.Errorf("failed to close: %v", err)
		return
	}
	log.Infof("Stream closed successfully for chunk id: %d\n", opts.chunkId)
}
