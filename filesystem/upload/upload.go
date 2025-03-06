package upload

import (
	"context"
	"os"
	"sync"

	"github.com/google/uuid"
	fm_pb "github.com/nanoDFS/client-sdk/filesystem/proto/master"
	"github.com/nanoDFS/client-sdk/utils/config"
	"google.golang.org/grpc"
)

type Uploader struct {
	filePath string
}

func NewUploader(filePath string) *Uploader {
	return &Uploader{filePath: filePath}
}

func (t *Uploader) Upload() (string, string, error) {
	info, _, err := t.openFile(t.filePath)
	if err != nil {
		return "", "", err
	}

	fileId, userId := uuid.NewString(), uuid.NewString()

	metaResponse, err := t.connectToMaster(fileId, userId, info)
	if err != nil {
		return "", "", err
	}

	t.streamMux(info, metaResponse)
	return fileId, userId, nil
}

func (t *Uploader) openFile(filePath string) (os.FileInfo, *os.File, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}
	file, _ := os.Open(filePath)
	return info, file, nil
}

func (t *Uploader) connectToMaster(fileId string, userId string, info os.FileInfo) (*fm_pb.UploadResp, error) {

	masterAddr := config.LoadConfig().Master.Addr
	conn, err := grpc.NewClient(masterAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := fm_pb.NewFileMetadataServiceClient(conn)
	resp, err := client.UploadFile(context.Background(), &fm_pb.FileUploadReq{FileId: fileId, UserId: userId, Size: info.Size()})
	return resp, err
}

func (t *Uploader) streamMux(fileInfo os.FileInfo, metaResponse *fm_pb.UploadResp) {
	chunkSize := config.LoadConfig().Chunk.Size
	chunkCount := fileInfo.Size() / chunkSize
	if fileInfo.Size()%chunkSize != 0 {
		chunkCount++
	}

	wg := &sync.WaitGroup{}

	for i := range chunkCount {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == chunkCount-1 {
			end = fileInfo.Size()
		}
		wg.Add(1)
		go func(start int64, end int64, i int64) {
			_, file, _ := t.openFile(t.filePath)
			t.stream(wg, file, start, end, metaResponse.ChunkServers[i].Address, int64(i), string(metaResponse.GetAccessToken()))
		}(start, end, i)
	}

	wg.Wait()
}
