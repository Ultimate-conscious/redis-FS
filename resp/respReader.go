package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) Readline() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) ReadInteger() (x int, n int, err error) {
	line, n, err := r.Readline()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {

	case ARRAY:
		return r.ReadArray()

	case BULK:
		return r.ReadBulk()

	default:
		fmt.Printf("Unknown type: %c\n", _type)
		return Value{}, nil
	}

}

func (r *Resp) ReadArray() (Value, error) {
	v := Value{}
	v.typ = "array"
	length, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Value, length)

	// for each line, parse and read the value
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		// append the value to the array
		v.array[i] = val
	}

	return v, nil

}

func (r *Resp) ReadBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	length, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	//reading the \r\n
	r.Readline()

	return v, nil

}
