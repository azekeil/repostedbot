package reposted

import (
	"fmt"
	"image"
	"net/http"

	"github.com/corona10/goimagehash"
)

func hashImageFromURL(url string) (*goimagehash.ImageHash, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got %d downloading %s", res.StatusCode, url)
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}
	return goimagehash.AverageHash(m)
}

func findRepost(hashMap ImgHashPost, hash *goimagehash.ImageHash, distance int) (*goimagehash.ImageHash, error) {
	for h := range hashMap.Iter() {
		loopHash := goimagehash.NewImageHash(h, goimagehash.AHash)
		d, err := hash.Distance(loopHash)
		if err != nil {
			return nil, err
		}
		if d <= distance {
			return loopHash, nil
		}
	}
	return nil, nil
}
