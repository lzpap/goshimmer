package transaction

import (
	"crypto/rand"
	"fmt"

	"github.com/mr-tron/base58"

	"github.com/iotaledger/goshimmer/packages/binary/marshalutil"
)

// Id is the data type that represents the identifier for a Transaction.
type Id [IdLength]byte

// IdFromBase58 creates an id from a base58 encoded string.
func IdFromBase58(base58String string) (id Id, err error) {
	// decode string
	bytes, err := base58.Decode(base58String)
	if err != nil {
		return
	}

	// sanitize input
	if len(bytes) != IdLength {
		err = fmt.Errorf("base58 encoded string does not match the length of a transaction id")

		return
	}

	// copy bytes to result
	copy(id[:], bytes)

	return
}

// IdFromBytes unmarshals an Id from a sequence of bytes.
func IdFromBytes(bytes []byte) (result Id, err error, consumedBytes int) {
	// parse the bytes
	marshalUtil := marshalutil.New(bytes)
	if idBytes, idErr := marshalUtil.ReadBytes(IdLength); idErr != nil {
		err = idErr

		return
	} else {
		copy(result[:], idBytes)
	}
	consumedBytes = marshalUtil.ReadOffset()

	return
}

// Parse is a wrapper for simplified unmarshaling of Ids from a byte stream using the marshalUtil package.
func ParseId(marshalUtil *marshalutil.MarshalUtil) (Id, error) {
	if id, err := marshalUtil.Parse(func(data []byte) (interface{}, error, int) { return IdFromBytes(data) }); err != nil {
		return Id{}, err
	} else {
		return id.(Id), nil
	}
}

// Random creates a random address, which can for example be used in unit tests.
func RandomId() (id Id) {
	// generate a random sequence of bytes
	idBytes := make([]byte, IdLength)
	if _, err := rand.Read(idBytes); err != nil {
		panic(err)
	}

	// copy the generated bytes into the result
	copy(id[:], idBytes)

	return
}

// Bytes marshals the Id into a sequence of bytes.
func (id Id) Bytes() []byte {
	return id[:]
}

// String creates a human readable version of the Id (for debug purposes).
func (id Id) String() string {
	return base58.Encode(id[:])
}

// IdLength contains the amount of bytes that a marshaled version of the Id contains.
const IdLength = 32