package netstrings

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

const (
	lengthDelim byte = ':'
	dataDelim   byte = ','
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) ReadNetstring() ([]byte, error) {
	length, err := r.r.ReadBytes(lengthDelim)
	if err != nil {
		return nil, err
	}
	l, err := strconv.Atoi(strings.TrimSuffix(string(length), string(lengthDelim)))
	if err != nil {
		return nil, err
	}
	ret := make([]byte, l)
	rd := ret
	for len(rd) > 0 {
		n, err := r.r.Read(rd)
		rd = rd[n:]
		if err != nil {
			return nil, err
		}
	}
	next, err := r.r.ReadByte()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if next != dataDelim {
		r.r.UnreadByte()
	}
	return ret, nil
}

func Decode(in []byte) ([][]byte, error) {
	rd := NewReader(bytes.NewReader(in))
	ret := make([][]byte, 0)
	for {
		d, err := rd.ReadNetstring()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		ret = append(ret, d)
	}
	return ret, nil
}

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteNetstring(b []byte) error {
	_, err := w.w.Write([]byte(strconv.Itoa(len(b))))
	if err != nil {
		return err
	}
	_, err = w.w.Write([]byte{lengthDelim})
	if err != nil {
		return err
	}
	_, err = w.w.Write(b)
	if err != nil {
		return err
	}
	_, err = w.w.Write([]byte{dataDelim})
	if err != nil {
		return err
	}
	return nil
}

func Encode(in ...[]byte) ([]byte, error) {
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	for _, d := range in {
		err := wr.WriteNetstring(d)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
