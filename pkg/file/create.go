package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// Write 指定したパスにファイルを作成する (Dirが存在しないなら合わせて作成する)
func Write(filePath string, content []byte) error {
	// ファイルのディレクトリパスを取得
	dirPath := filepath.Dir(filePath)

	// ディレクトリが存在しなければ作成
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
