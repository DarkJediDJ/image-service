package img

import (
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/DarkJediDJ/image-service/processing-service/internal/cloud"
	"github.com/nfnt/resize"
)

func Resize(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	file.Close()

	var images []string

	for i := 100; i < 1200; i += 500 {
		img := resize.Resize(uint(i), 0, img, resize.Lanczos3)
		out, err := ioutil.TempFile(cloud.FilePath, "img-*.jpeg")
		if err != nil {
			return nil, err
		}
		defer out.Close()
		if err = jpeg.Encode(out, img, nil); err != nil {
			return nil, err
		}
		images = append(images, out.Name())
	}

	return images, err
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
