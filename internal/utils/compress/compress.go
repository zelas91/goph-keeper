package compress

import (
	"compress/gzip"
	"github.com/zelas91/goph-keeper/internal/logger"
	"sync"
)

type Compress struct {
	compress *sync.Pool
	log      logger.Logger
}

func NewCompress(log logger.Logger) *Compress {
	return &Compress{compress: &sync.Pool{New: func() any {
		writer, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
		if err != nil {
			return gzip.NewWriter(nil)
		}
		return writer
	}},
		log: log}
}

func (c *Compress) Writer() *gzip.Writer {
	if v := c.compress.Get(); v != nil {
		return v.(*gzip.Writer)
	}
	writer, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		c.log.Errorf("Failed to create gzip writer err: %v", err)
		return nil
	}

	return writer
}

func (c *Compress) Release(writer *gzip.Writer) {
	defer func(w *gzip.Writer) {
		if err := writer.Close(); err != nil {
			c.log.Errorf("Failed to close gzip writer:", err)
			return
		}
	}(writer)
	c.compress.Put(writer)
}

type Decompress struct {
	decompress *sync.Pool
}

func NewDecompress() *Decompress {
	return &Decompress{decompress: &sync.Pool{}}
}

func (d *Decompress) Reader() *gzip.Reader {
	reader := d.decompress.Get()
	if reader == nil {
		return &gzip.Reader{}
	}

	return reader.(*gzip.Reader)
}

func (d *Decompress) Release(reader *gzip.Reader) {
	d.decompress.Put(reader)
}
