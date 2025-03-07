package delete

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"

	fs_pb "github.com/nanoDFS/client-sdk/filesystem/proto/chunkserver"
	fm_pb "github.com/nanoDFS/client-sdk/filesystem/proto/master"
	"github.com/nanoDFS/client-sdk/utils/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Deletor struct {
	fileId string
	userId string
}

func NewDeletor(fileId string, userId string) *Deletor {
	return &Deletor{
		fileId,
		userId,
	}
}
func (t *Deletor) Delete() error {
	metaResponse, err := t.connectToMaster()
	if err != nil {
		log.Fatalf("failed to download: %v", err)
	}
	t.deleteReqToCS(metaResponse)
	return nil
}

func (t *Deletor) connectToMaster() (*fm_pb.DeleteResp, error) {
	masterAddr := config.LoadConfig().Master.Addr
	conn, err := grpc.NewClient(masterAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := fm_pb.NewFileMetadataServiceClient(conn)
	resp, err := client.DeleteFile(context.Background(), &fm_pb.FileDeleteReq{FileId: t.fileId, UserId: t.userId})
	return resp, err
}

func (t *Deletor) deleteReqToCS(metadata *fm_pb.DeleteResp) {
	chunkCount := len(metadata.ChunkServers)

	wg := &sync.WaitGroup{}
	for i := range chunkCount {
		wg.Add(1)
		go t.delete(wg, metadata.ChunkServers[i].Address, int64(i), string(metadata.GetAccessToken()))
	}

	wg.Wait()
}

func (t *Deletor) delete(wg *sync.WaitGroup, address string, chunkId int64, token string) {
	defer wg.Done()

	conn, err := grpc.NewClient(address, grpc.WithInsecure())
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
	_, err = client.Delete(ctx, &fs_pb.DeleteReq{})
	if err != nil {
		log.Errorf("failed to delete file %s : %v", t.fileId, err)
		return
	}
}
