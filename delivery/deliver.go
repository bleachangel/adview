package delivery

import (
	"adview/resource"
	"crypto/md5"
	"encoding/hex"
)

type Delivery struct {
	res resource.Storage //资源存储
	db  *sql.DB
}

func (d *Delivery) upload(w http.ResponseWriter, r *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	ct := r.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil || mt != "multipart/form-data" {
		http.Error(w, "Unsupported Media Type", 415)
		return
	}

	status := "{\"res\": %d, \"msg\":\"%s\"}"
	mr, err := r.MultipartReader()
	if err != nil {
		glog.Errorf("create multipart reader error: %v", err)
		w.Write([]byte(fmt.Sprintf(status, 503, err.Error())))
		return
	}
	form, err := mr.ReadForm(1024)
	if err != nil {
		glog.Errorf("read multipart error: %v", err)
		w.Write([]byte(fmt.Sprintf(status, 503, err.Error())))
		return
	}

	//step 1 保存文件
	var success int
	var lastErr error
	for _, files := range form.File {
		for _, file := range files {
			rsFile, err := file.Open()
			if err != nil {
				glog.Warningf("can't open resource file %s: %v", file.Filename, err)
				lastErr = err
				continue
			}

			size, url, err := d.res.Save(rsFile)
			if err != nil {
				glog.Warningf("save resource file %s error: %v", file.Filename, err)
				lastErr = err
				continue
			}
			err = d.deliver(form, size, url)
			if err != nil {
				lastErr = err
				continue
			}
			success++
		}
	}

	if success == 0 {
		w.Write([]byte(fmt.Sprintf(status, 503, lastErr.Error())))
	} else {
		w.Write([]byte(fmt.Sprintf(status, 200, "OK")))
	}
}

func (d *Delivery) deliver(form *mime.Form, size int64, url string) error {
	adMeta, err := d.parseForm(form)
	if err != nil {
		return err
	}

	adMeta.Size = size
	adMeta.Url = url
	adMeta.CreateTime = time.Now()
	adMeta.AdToken = d.generateToken(url, size)

	//TODO 写入数据库
	//TODO 调用推送接口
}

func (d *Delivery) generateToken(url string, size int64) string {
	hash = md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", url, size)))
	return hex.EncodeToString(hash.Sum(nil))
}

func (d *Delivery) parseForm(form *mime.Form) (*types.AdMeta, error) {
	var (
		ok     bool
		values []string
	)
	err := errors.New("Invalid post parameter")
	adMeta := &types.AdMeta{}

	ok, values = form.Value["name"]
	if !ok || len(values) == 0 {
		return nil, err
	}
	adMeta.Name = values[0]

	ok, values = form.Value["userid"]
	if !ok || len(values) == 0 {
		return nil, err
	}
	adMeta.UserId = values[0]

	ok, values = form.Value["at"]
	if !ok || len(values) == 0 {
		return nil, err
	}
	adMeta.At, _ = strconv.Atoi(values[0])
	if adMeta.At != 0 {
		return nil, err
	}

	ok, values = form.Value["aw"]
	if !ok || len(values) == 0 {
		return nil, err
	}
	adMeta.Aw, _ = strconv.Atoi(values[0])

	ok, values = form.Value["ah"]
	if !ok || len(values) == 0 {
		return nil, err
	}
	adMeta.Ah, _ = strconv.Atoi(values[0])

	ok, values = form.Value["priority"]
	if !ok || len(values) == 0 {
		adMeta.Priority = 1
	} else {
		adMeta.Priority, _ = strconv.Atoi(values[0])
	}

	return adMeta, nil
}

func Run(res resource.Storage) {
	d := &Delivery{res: res}
	http.HandleFunc("/v1/adview/deliver", d.upload)
}
