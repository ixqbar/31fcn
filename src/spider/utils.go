package spider

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"path/filepath"
	"github.com/Pallinder/go-randomdata"
)

func CheckFileIsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if fi.IsDir() == false {
		return false, errors.New("target file is not folder")
	}

	return true, nil
}

func GetFileSize(file string) (int64, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	if fi.IsDir() {
		return 0, errors.New(fmt.Sprintf("target file %s is not file", file))
	}

	return fi.Size(), nil
}

func InStringArray(value string, arrays []string) bool {
	for _, v := range arrays {
		if v == value {
			return true
		}
	}

	return false
}

func GetFileMD5sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func HasIntersection(a []string, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	t := strings.Join(b, "%") + "%"
	for _, v := range a {
		if strings.Contains(t, v+"%") {
			return true
		}
	}

	return false
}

func IsFalse(needle string) bool {
	if len(needle) == 0 {
		return true
	}

	haystack := []interface{}{
		false,
		0,
		"false",
		"",
	}

	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func Rc4Decrypt(content []byte, key []byte) ([]byte, error) {
	rc4Cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(content))
	rc4Cipher.XORKeyStream(plainText, content)

	return plainText, nil
}

func Left(str string, length int, pad string) string {
	return strings.Repeat(pad, length-len(str)) + str
}

func Right(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-len(str))
}

func GetTraceRealUrl(refererUrl, url string) string {
	realUrl := ""
	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 10 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			realUrl = req.URL.String()
			return errors.New("stopped redirects")
		},
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Print(err)
		return realUrl
	}

	request.Header.Set("Referer", refererUrl)
	request.Header.Set("user-agent", randomdata.UserAgentString())

	response, err := httpClient.Do(request)
	if err != nil {
		Logger.Print(err)
		return realUrl
	}

	if response.StatusCode == 302 {
		realUrl = response.Request.URL.String()
	}

	return realUrl
}

func ChownR(path string, uid, gid int) (error) {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = Chown(name, uid, gid)
		} else {
			Logger.Printf("walk %s failed %s", name, err)
		}
		return err
	})
}


func Chown(name string ,uid, gid int) (error) {
	err := os.Chown(name, uid, gid)
	if err != nil {
		Logger.Printf("chown %s failed", name, err)
	}

	return err
}

func UrlLastPath(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
