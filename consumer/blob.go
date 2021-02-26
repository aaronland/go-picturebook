package consumer

import (
	"context"
	"github.com/aaronland/go-mimetypes"
	"github.com/aaronland/go-picturebook"
	"github.com/aaronland/go-picturebook/picture"
	"gocloud.dev/blob"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type BlobConsumer struct {
	Consumer
	bucket *blob.Bucket
	mu     *sync.RWMutex
}

func init() {

	ctx := context.Background()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		err := RegisterConsumer(ctx, scheme, NewBlobConsumer)

		if err != nil {
			panic(err)
		}
	}
}

func NewBlobConsumer(ctx context.Context, uri string) (Consumer, error) {

	bucket, err := blob.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	mu := new(sync.RWMutex)

	c := &BlobConsumer{
		bucket: bucket,
		mu:     mu,
	}

	return c, nil
}

func (c *BlobConsumer) GatherPictures(ctx context.Context, pb_opts *picturebook.PictureBookOptions, uris ...string) ([]*picture.PictureBookPicture, error) {

	pictures := make([]*picture.PictureBookPicture, 0)

	var list func(context.Context, *blob.Bucket, string) error

	process_file := func(ctx context.Context, b *blob.Bucket, path string) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		// START OF sudo put me in a generic method that works with an io.ReadSeeker (the image)
		// and a path...
		
		abs_path := path

		is_image := false

		ext := filepath.Ext(abs_path)
		ext = strings.ToLower(ext)

		for _, t := range mimetypes.TypesByExtension(ext) {
			if strings.HasPrefix(t, "image/") {
				is_image = true
				break
			}
		}

		if !is_image {

			if pb_opts.Verbose {
				log.Printf("%s (%s) does not appear to be an image, skipping\n", abs_path, ext)
			}

			return nil
		}

		if pb_opts.Filter != nil {

			// WHAT WHAT SOURCE VS CONSUMER
			ok, err := pb_opts.Filter.Continue(ctx, pb_opts.Source, abs_path)

			if err != nil {
				log.Printf("Failed to filter %s, %v\n", abs_path, err)
				return nil
			}

			if !ok {
				return nil
			}

			if pb_opts.Verbose {
				log.Printf("Include %s\n", abs_path)
			}
		}

		caption := ""

		if pb_opts.Caption != nil {

			// WHAT WHAT SOURCE VS CONSUMER
			txt, err := pb_opts.Caption.Text(ctx, pb_opts.Source, abs_path)

			if err != nil {
				log.Printf("Failed to generate caption text for %s, %v\n", abs_path, err)
				return nil
			}

			caption = txt
		}

		var final_bucket *blob.Bucket
		final_path := abs_path

		if pb_opts.PreProcess != nil {

			if pb_opts.Verbose {
				log.Printf("Processing %s\n", abs_path)
			}

			// WHAT WHAT SOURCE VS CONSUMER
			processed_path, err := pb_opts.PreProcess.Transform(ctx, pb_opts.Source, pb_opts.Temporary, abs_path)

			if err != nil {
				log.Printf("Failed to process %s, %v\n", abs_path, err)
				return nil
			}

			if pb_opts.Verbose {
				log.Printf("After processing %s becomes %s\n", abs_path, processed_path)
			}

			if processed_path != "" && processed_path != abs_path {
				// pb.tmpfiles = append(pb.tmpfiles, processed_path)
				final_path = processed_path
				final_bucket = pb_opts.Temporary
			}
		}

		// END OF sudo put me in a generic method that works with an io.ReadSeeker (the image)
		
		c.mu.Lock()
		defer c.mu.Unlock()

		if pb_opts.Verbose {
			log.Printf("Append %s (%s) to list for processing\n", final_path, abs_path)
		}

		// WHAT WHAT BUCKETS AT ALL...

		pic := &picture.PictureBookPicture{
			Source:  abs_path,
			Bucket:  final_bucket,
			Path:    final_path,
			Caption: caption,
		}

		pictures = append(pictures, pic)
		return nil
	}

	list = func(ctx context.Context, bucket *blob.Bucket, prefix string) error {

		iter := bucket.List(&blob.ListOptions{
			Delimiter: "/",
			Prefix:    prefix,
		})

		for {
			obj, err := iter.Next(ctx)

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			path := obj.Key

			if obj.IsDir {

				err := list(ctx, bucket, path)

				if err != nil {
					return err
				}

				continue
			}

			err = process_file(ctx, bucket, path)

			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, path := range uris {

		err := list(ctx, pb_opts.Source, path)

		if err != nil {
			return nil, err
		}
	}

	return pictures, nil
}
