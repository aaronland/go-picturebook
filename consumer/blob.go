package consumer

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
)

type BlobConsumer struct {
	Consumer
}

func init() {

	ctx := context.Background()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		err := consumer.RegisterConsumer(ctx, scheme, NewBlobConsumer)

		if err != nil {
			panic(err)
		}
	}
}

func NewBlobConsumer(ctx context.Context, uri string) (Consumer, error) {

	return nil, fmt.Errorf("Not implemented")
}

func (c *BlobConsumer) GatherPictures(ctx context.Context, uris ...string) ([]*picture.PictureBookPicture, error) {

	return nil, fmt.Errorf("Not implemented")

	/*
		pictures := make([]*picture.PictureBookPicture, 0)

		var list func(context.Context, *blob.Bucket, string) error

		process_file := func(ctx context.Context, b *blob.Bucket, path string) error {

			select {
			case <-ctx.Done():
				return nil
			default:
				// pass
			}

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

				if pb.Options.Verbose {
					log.Printf("%s (%s) does not appear to be an image, skipping\n", abs_path, ext)
				}

				return nil
			}

			if pb.Options.Filter != nil {

				ok, err := pb.Options.Filter.Continue(ctx, pb.Options.Source, abs_path)

				if err != nil {
					log.Printf("Failed to filter %s, %v\n", abs_path, err)
					return nil
				}

				if !ok {
					return nil
				}

				if pb.Options.Verbose {
					log.Printf("Include %s\n", abs_path)
				}
			}

			caption := ""

			if pb.Options.Caption != nil {

				txt, err := pb.Options.Caption.Text(ctx, pb.Options.Source, abs_path)

				if err != nil {
					log.Printf("Failed to generate caption text for %s, %v\n", abs_path, err)
					return nil
				}

				caption = txt
			}

			var final_bucket *blob.Bucket
			final_path := abs_path

			if pb.Options.PreProcess != nil {

				if pb.Options.Verbose {
					log.Printf("Processing %s\n", abs_path)
				}

				processed_path, err := pb.Options.PreProcess.Transform(ctx, pb.Options.Source, pb.Options.Temporary, abs_path)

				if err != nil {
					log.Printf("Failed to process %s, %v\n", abs_path, err)
					return nil
				}

				if pb.Options.Verbose {
					log.Printf("After processing %s becomes %s\n", abs_path, processed_path)
				}

				if processed_path != "" && processed_path != abs_path {
					pb.tmpfiles = append(pb.tmpfiles, processed_path)
					final_path = processed_path
					final_bucket = pb.Options.Temporary
				}
			}

			pb.Mutex.Lock()
			defer pb.Mutex.Unlock()

			if pb.Options.Verbose {
				log.Printf("Append %s (%s) to list for processing\n", final_path, abs_path)
			}

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

			err := list(ctx, pb.Options.Source, path)

			if err != nil {
				return nil, err
			}
		}

		return pictures, nil
	*/
}
