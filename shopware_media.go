package main

import (
    "crypto/md5"
    "fmt"
    "github.com/pkg/errors"
    "github.com/sirupsen/logrus"
    "path"
    "regexp"
    "strconv"
    "strings"
)

const (
    MediaImage = "/media/image/"
)

var thumbnailPathRegex = regexp.MustCompile(`\w{2}\/\w{2}\/\w{2}\/([\w-]+)_(\d+)x(\d+)(@2x)?\.(\w+)`)
var formats = []string{"JPG", "PNG"}

func replaceBadWords(input string) string {
    return strings.ReplaceAll(input, "/ad/", "/g0/")
}

type ThumbnailPath struct {
    retina   bool
    fileName string
    format   string
    height   int
    width    int
}

var ErrNoThumbnail = errors.New("not a thumbnail path")

func ParseThumbnailPath(path string) (err error, tp *ThumbnailPath) {
    if strings.HasPrefix(path, MediaImage) {
        path = strings.Replace(path, MediaImage, "", 1)
    }

    logrus.Info(path)

    var regexResult = thumbnailPathRegex.FindStringSubmatch(path)

    if regexResult == nil {
        return ErrNoThumbnail, nil
    }
    tp = &ThumbnailPath{}

    if regexResult[4] == "@2x" {
        tp.retina = true
    }

    tp.format = strings.ToUpper(regexResult[5])
    if !StringArray(formats).Contains(tp.format) {
        tp.format = formats[0]
    }

    tp.fileName = regexResult[1] + "." + strings.ToLower(tp.format)

    tp.height, err = strconv.Atoi(regexResult[2])
    if err != nil {
        return err, nil
    }

    tp.width, err = strconv.Atoi(regexResult[3])
    if err != nil {
        return err, nil
    }

    return nil, tp
}

func (tp *ThumbnailPath) Encode() string {
    var md5String = fmt.Sprintf("%x", md5.Sum([]byte(MediaImage+tp.fileName)))

    url := path.Join(md5String[0:2], md5String[2:4], md5String[4:6], tp.fileName)

    return replaceBadWords(url)
}
