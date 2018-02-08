package resource

import (
	"io"
)

type Storage interface {
	Save(src io.Reader) (int64, string, error)
}

type localStorage struct {
	basePath string
}

func NewStorage() Storage {
	store := &localStorage{}
	return store
}

func (store *localStorage) Save(src io.Reader) (int64, string, error) {
	h := md5.New()
	if _, err := io.Copy(h, src); err != nil {
		return 0, "", err
	}

	//生成新的文件名
	dstFileName := hex.EncodeToString(h.Sum(nil))
	now := time.Now()
}
