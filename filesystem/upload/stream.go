package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/nanoDFS/client-sdk/crypto"
	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type streamOpts struct {
	wg         *sync.WaitGroup
	file       *os.File
	start      int64
	end        int64
	serverAddr string
	chunkId    int64
	token      string
	key        crypto.CryptoKey
}

func (t *Uploader) stream(opts streamOpts) {
	defer opts.wg.Done()
	conn, err := grpc.NewClient(opts.serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		log.Errorf("%v", err)
		return
	}
	defer conn.Close()
	md := metadata.Pairs(
		"auth", string(opts.token),
		"chunk_id", fmt.Sprintf("%d", opts.chunkId),
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := fs_pb.NewFileStreamingServiceClient(conn)
	stream, err := client.Write(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	if _, err := opts.file.Seek(opts.start, io.SeekStart); err != nil {
		log.Errorf("%v", err)
		return
	}
	buff := make([]byte, 1024)
	for {
		n, err := opts.file.Read(buff)
		if err != nil && err.Error() != "EOF" {
			log.Errorf("%v", err)
			return
		}
		currentPos, _ := opts.file.Seek(0, io.SeekCurrent)
		if n == 0 || currentPos > opts.end {
			break
		}
		encryptedData, err := crypto.NewEncryptor().Encrypt(buff[:n], t.key.Key, t.key.Nonce)
		if err != nil {
			log.Errorf("failed to encrypt: %v", err)
			break
		}
		if err := stream.Send(&fs_pb.Payload{Data: encryptedData}); err != nil {
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
