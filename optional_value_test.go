package gopt_mongo

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"testing"
)

type myStruct struct {
	Foo Optional[string] `json:"foo" bson:"foo"`
	Bar Optional[int64]  `json:"bar" bson:"bar"`
}

func TestInitialises(t *testing.T) {
	my := &myStruct{}
	require.NotNil(t, my.Foo)
	require.NotNil(t, my.Foo.Optional)
}

func TestOptional_UnmarshalBSONValue(t *testing.T) {
	data := []byte{7, 0, 0, 0, 'm', 'a', 'r', 'r', 'o', 'w', 0}
	my := &myStruct{}
	err := my.Foo.UnmarshalBSONValue(bsontype.String, data)
	require.NoError(t, err)
	require.True(t, my.Foo.IsPresent())
	require.True(t, my.Foo.WasSet())
	v, ok := my.Foo.GetOk()
	require.True(t, ok)
	require.Equal(t, "marrow", v)

	data = nil
	err = my.Foo.UnmarshalBSONValue(bsontype.String, data)
	require.Error(t, err)

	data = []byte{7, 0, 0, 0, 'm', 'a', 'r', 'r', 'o', 'w'}
	err = my.Foo.UnmarshalBSONValue(bsontype.String, data)
	require.Error(t, err)

	data = []byte{4, 0, 0, 0, 0, 0, 0, 0, 0}
	err = my.Bar.UnmarshalBSONValue(bsontype.Int64, data)
	require.NoError(t, err)

	err = my.Foo.UnmarshalBSONValue(bsontype.Int64, data)
	require.Error(t, err)

	err = my.Foo.unmarshalBSONValue(func() {})
	require.Error(t, err)
}

func TestOptional_MarshalBSONValue(t *testing.T) {
	my := &myStruct{}
	ty, by, err := my.Foo.MarshalBSONValue()
	require.NoError(t, err)
	require.Equal(t, bsontype.Null, ty)
	require.Nil(t, by)

	my.Foo.OrElseSet("abc")
	ty, by, err = my.Foo.MarshalBSONValue()
	require.NoError(t, err)
	require.Equal(t, bsontype.String, ty)
	require.Equal(t, 8, len(by))
	require.Equal(t, []byte{4, 0, 0, 0, 'a', 'b', 'c', 0}, by)
}

func TestOptional_Unmarshal(t *testing.T) {
	// check that underlying Optional is still called on, for example, Unmarshal
	my := &myStruct{}
	err := json.Unmarshal([]byte(`{"foo":"marrow","bar":10}`), my)
	require.NoError(t, err)
	require.True(t, my.Foo.IsPresent())
}
