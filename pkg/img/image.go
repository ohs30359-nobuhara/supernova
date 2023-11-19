package img

import (
	"bytes"
	"image"
	"image/png"
	"os"
)

// ByteArrayToImage Byteを画像に変換する
func ByteArrayToImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImageToFile imageをファイルとして保存する
func SaveImageToFile(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
