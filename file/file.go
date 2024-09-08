package file

import (
	"os"
	"path"
)

func WriteToFile(filepath string, content string) error {
	if err := os.MkdirAll(path.Dir(filepath), 0o700); err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	return err
}
