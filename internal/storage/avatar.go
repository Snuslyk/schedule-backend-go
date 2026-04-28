package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"time"

	"github.com/disintegration/imaging"
)

type CropData struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Size int `json:"size"`
}

type AvatarService struct {
	storage Storage
}

func NewAvatarService(storage Storage) *AvatarService {
	return &AvatarService{storage: storage}
}

func (s *AvatarService) Upload(ctx context.Context, userID string, fileData []byte, crop CropData) error {
	img, _, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	if crop.Size <= 0 {
		crop.Size = bounds.Dx()
		if bounds.Dy() < crop.Size {
			crop.Size = bounds.Dy()
		}
	}
	if crop.X < 0 {
		crop.X = 0
	}
	if crop.Y < 0 {
		crop.Y = 0
	}
	if crop.X+crop.Size > bounds.Dx() {
		crop.Size = bounds.Dx() - crop.X
	}
	if crop.Y+crop.Size > bounds.Dy() {
		crop.Size = bounds.Dy() - crop.Y
	}

	rect := image.Rect(crop.X, crop.Y, crop.X+crop.Size, crop.Y+crop.Size)
	cropped := imaging.Crop(img, rect)
	resized := imaging.Resize(cropped, 512, 512, imaging.Lanczos)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("failed to encode jpeg: %w", err)
	}

	key := fmt.Sprintf("avatars/%s.jpg", userID)
	return s.storage.Upload(ctx, key, buf.Bytes(), "image/jpeg")
}

func (s *AvatarService) GetAvatarURL(ctx context.Context, userID string) (*string, error) {
	key := fmt.Sprintf("avatars/%s.jpg", userID)
	exists, err := s.storage.Exists(ctx, key)
	if err != nil || !exists {
		return nil, nil
	}

	url, err := s.storage.GetPresignedURL(ctx, key, 1*time.Hour)
	if err != nil {
		return nil, err
	}
	return &url, nil
}
