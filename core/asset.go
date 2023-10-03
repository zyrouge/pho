package core

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/zyrouge/pho/utils"
)

type AssetDownloadFunc func() (io.ReadCloser, error)

type Asset struct {
	Source   string
	Size     int64
	Download AssetDownloadFunc
}

func NetworkAssetDownload(url string) AssetDownloadFunc {
	return func() (io.ReadCloser, error) {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		return res.Body, nil
	}
}

func LocalAssetDownload(name string) AssetDownloadFunc {
	return func() (io.ReadCloser, error) {
		file, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

type NetworkAssetMetadata struct {
	Size int64
}

func ExtractNetworkAssetMetadata(url string) (*NetworkAssetMetadata, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	metadata := &NetworkAssetMetadata{
		Size: res.ContentLength,
	}
	return metadata, nil
}

func ChooseAptAppImageAsset[T any](assets []T, assetNameFunc func(*T) string) (int, *T) {
	arch := utils.GetSystemArch()
	var fallback *T
	for i := range assets {
		asset := &assets[i]
		name := strings.ToLower(assetNameFunc(asset))
		if !strings.HasSuffix(name, ".appimage") {
			continue
		}
		matchedArch := extractArch(name)
		if matchedArch == arch {
			return 2, asset
		}
		// no arch probably means they didnt include it
		if matchedArch == "" {
			fallback = asset
		}
	}
	if fallback != nil {
		return 1, fallback
	}
	return 0, nil
}

func extractArch(name string) string {
	for arch, aliases := range utils.ArchMap {
		for _, x := range aliases {
			if strings.Contains(name, x) {
				return arch
			}
		}
	}
	return ""
}
