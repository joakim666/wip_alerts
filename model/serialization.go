package model

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/golang/glog"
)

func serialize(obj interface{}) ([]byte, error) {
	glog.Infof("Serialization")
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize object: %s", err)
	}

	return buf.Bytes(), nil
}

// copy the bytes and deserialize the object
func deserialize(src *[]byte, obj interface{}) error {
	glog.Infof("Deserialization")

	// make a copy of the bytes as this comes from bolt and the object will be used outside
	// of the transaction
	b := make([]byte, len(*src))
	copy(b, *src)

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(obj)
	if err != nil {
		glog.Infof("Failed to deserialize object: %s", err)
		return err
	}

	return nil
}
