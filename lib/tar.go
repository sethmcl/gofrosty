package lib

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

// ExtractTar extracts tar archive
func ExtractTar(f io.Reader, target string, truncate int) error {
	ctx := GetContext()

	if IsFile(target) {
		return fmt.Errorf("Target exists and is not a directory (%s)", target)
	}

	if !IsDir(target) {
		err := os.MkdirAll(target, os.ModePerm)
		if err != nil {
			return cleanup(target, err)
		}
	}

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return cleanup(target, err)
	}
	defer gzf.Close()

	tarReader := tar.NewReader(gzf)

	i := 0
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return cleanup(target, err)
		}

		name, err := TruncatePath(header.Name, truncate)
		if err != nil {
			return err
		}
		absPath := path.Join(target, name)

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(absPath, os.ModePerm)
			if err != nil {
				ctx.Debug("Unable to create %s with mode %s", absPath, header.Mode)
				return cleanup(target, err)
			}
			continue
		case tar.TypeReg:
			os.MkdirAll(path.Dir(absPath), os.ModePerm)
			writer, err := os.Create(absPath)
			if err != nil {
				return cleanup(target, err)
			}

			io.Copy(writer, tarReader)
			writer.Close()
		default:
			err := fmt.Errorf("untar failed type: %c in file %s", header.Typeflag, absPath)
			return cleanup(target, err)
		}

		i++
	}

	return nil
}

func cleanup(target string, err error) error {
	os.RemoveAll(target)
	return err
}
