package file

import "os"

func FileOverwrite(path string, content string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}
