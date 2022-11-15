package gopt_mongo

import (
	"encoding/json"
	"github.com/go-andiamo/gopt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Optional[T any] struct {
	gopt.Optional[T]
}

func (v *Optional[T]) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	vr := bsonrw.NewBSONValueReader(btype, data)
	dec, _ := bson.NewDecoder(vr)
	var dv interface{}
	err := dec.Decode(&dv)
	if err != nil {
		return err
	}
	return v.unmarshalBSONValue(dv)
}

func (v *Optional[T]) unmarshalBSONValue(iv interface{}) error {
	jdata, err := json.Marshal(&iv)
	if err != nil {
		return err
	}
	return v.UnmarshalJSON(jdata)
}

func (v Optional[T]) MarshalBSONValue() (bsontype.Type, []byte, error) {
	ov, ok := v.GetOk()
	if !ok {
		return bsontype.Null, nil, nil
	}
	return bson.MarshalValue(ov)
}
