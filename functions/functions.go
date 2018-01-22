package functions

type PictureBookFilterFunc func(string) (bool, error)

type PictureBookPreProcessFunc func(string) (string, error)

type PictureBookCaptionFunc func(string) (string, error)
