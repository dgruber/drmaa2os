package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/images"
)

func ListContainerImages(ctx context.Context) ([]string, error) {
	imageSummary, err := images.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	images := make([]string, 0, len(imageSummary))
	for _, is := range imageSummary {
		images = append(images, is.Names...)
	}
	return images, nil
}
