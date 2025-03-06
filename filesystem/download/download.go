package download

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	fm_pb "github.com/nanoDFS/client-sdk/filesystem/proto/master"
	"github.com/nanoDFS/client-sdk/utils/config"
	"google.golang.org/grpc"
)

type Downloader struct {
	fileId   string
	userId   string
	filePath string
}

func NewDownloader(fileId string, userId string, filePath string) *Downloader {
	return &Downloader{
		fileId:   fileId,
		userId:   userId,
		filePath: filePath,
	}
}

func (t *Downloader) Download() error {
	metaResponse, err := t.connectToMaster()
	if err != nil {
		log.Fatalf("failed to download: %v", err)
	}

	t.streamMux(metaResponse)
	return nil
}

func (t *Downloader) connectToMaster() (*fm_pb.DownloadResp, error) {
	masterAddr := config.LoadConfig().Master.Addr
	conn, err := grpc.NewClient(masterAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := fm_pb.NewFileMetadataServiceClient(conn)
	resp, err := client.DownloadFile(context.Background(), &fm_pb.FileDownloadReq{FileId: t.fileId, UserId: t.userId})
	if err != nil {
		fmt.Println(err)
	}
	return resp, err
}

func (t *Downloader) streamMux(metaResponse *fm_pb.DownloadResp) {
	chunkCount := len(metaResponse.ChunkServers)

	wg := &sync.WaitGroup{}
	for i := range chunkCount {
		wg.Add(1)
		go t.stream(wg, metaResponse.ChunkServers[i].Address, int64(i), string(metaResponse.GetAccessToken()))
	}

	wg.Wait()
}
