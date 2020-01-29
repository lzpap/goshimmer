package drng

import (
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/drand/drand/core"
	"github.com/drand/drand/protobuf/drand"
	"github.com/golang/protobuf/proto"
	"github.com/iotaledger/goshimmer/packages/model/value_transaction"
	"github.com/iotaledger/hive.go/typeutils"
	"github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/stretchr/testify/require"
)

func TestValidDRNG(t *testing.T) {
	signature := []byte("a85cc3216189fe520a2bf83c3e369b9d2c4fda3009e65e7fa2f17cc2d6c3f015dd9aaddc0d9117466d939d0827811226046d76c2a8f62e392d4cc50dbe3d5c915041a4939118ee35224aba2bb9bf974e6234d2d81c106559d79053cbb1b7f10d")
	h := core.RandomnessHash()
	h.Write(signature)
	randomness := h.Sum(nil)
	m := &drand.PublicRandResponse{
		Previous:   signature,
		Round:      1,
		Signature:  signature,
		Randomness: randomness,
	}
	data, err := proto.Marshal(m)
	require.NoError(t, err)

	q := hex.EncodeToString(data)

	size := strconv.Itoa((len(q)))

	buffer := make([]byte, 2187)
	copy(buffer, typeutils.StringToBytes(size+q))

	trytes, err := trinary.BytesToTrytes(buffer)
	require.NoError(t, err)

	tx := value_transaction.New()
	tx.SetHead(true)
	tx.SetTail(true)
	err = address.ValidAddress(defaultAddress)
	require.NoError(t, err)

	tx.SetAddress(defaultAddress)
	tx.SetSignatureMessageFragment(trytes)
	tx.SetValue(0)
	tx.SetTimestamp(uint(time.Now().Unix()))

	var s []byte
	s, err = hasValidData(tx)
	require.NoError(t, err)

	pb := &drand.PublicRandResponse{}
	err = proto.Unmarshal(s, pb)
	require.NoError(t, err)
}
