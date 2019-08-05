package main

import (
    "github.com/disintegration/imaging"
    "github.com/sirupsen/logrus"
    "io"
    "os"
    "path"
)

func CreateThumbnail(thumbnailPath *ThumbnailPath, writer io.Writer) error {
    file, err := os.Open(path.Join(srcFolder, thumbnailPath.fileName))
    if err != nil {
        return err
    }

    image, err := imaging.Decode(file)
    if err != nil {
        return err
    }

    logrus.Info(thumbnailPath.format)
    extension, err := imaging.FormatFromExtension(thumbnailPath.format)
    if err != nil {
        return err
    }

    height, width := calculateImageSize(image.Bounds().Max.Y, image.Bounds().Max.X, thumbnailPath.height, thumbnailPath.width)

    resize := imaging.Resize(image, width, height, imaging.Linear)



    err = imaging.Encode(writer, resize, extension)
    if err != nil {
        return err
    }

    return nil
}

func calculateImageSize(originalHeight, originalWidth, newHeight, newWidth int) (int, int) {
    floatOriginalHeight := float64(originalHeight)
    floatOriginalWidth := float64(originalWidth)
    floatNewHeight := float64(newHeight)
    floatNewWidth := float64(newWidth)

    ratio := floatOriginalWidth / floatOriginalHeight
    if floatNewWidth/floatNewHeight > ratio {
        floatNewWidth = floatNewHeight * ratio
    } else {
        floatNewHeight = floatNewWidth / ratio
    }

    return int(floatNewHeight), int(floatNewWidth)
}
