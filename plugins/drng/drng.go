package drng

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/drand/drand/protobuf/drand"
	"github.com/golang/protobuf/proto"
	"github.com/iotaledger/goshimmer/packages/model/value_transaction"
	"github.com/iotaledger/goshimmer/plugins/tangle"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/typeutils"
	"github.com/iotaledger/iota.go/trinary"
)

const (
	defaultAddress = "RANDOM99NUMBER9999999999999999999999999999999999999999999999999999999999999999999"
)

var (
	last *[]byte
)

func configureDRNG() {
	tangle.Events.TransactionSolid.Attach(events.NewClosure(func(tx *value_transaction.ValueTransaction) {
		var d []byte
		var err error
		if hasValidAddress(tx) != nil {
			return
		}
		if d, err = hasValidData(tx); err != nil {
			return
		}
		pb := &drand.PublicRandResponse{}
		err = proto.Unmarshal(d, pb)
		if err != nil {
			return
		}
		log.Info("New Random:", hex.EncodeToString(pb.GetRandomness()))
	}))
}

func hasValidAddress(tx *value_transaction.ValueTransaction) error {
	if tx.GetAddress() != defaultAddress {
		return errors.New("Not default DRNG address")
	}
	return nil
}

func hasValidData(tx *value_transaction.ValueTransaction) ([]byte, error) {
	buf, err := trinary.TrytesToBytes(tx.GetSignatureMessageFragment())
	if err != nil {
		return nil, err
	}

	f := typeutils.BytesToString(buf)
	l, err := strconv.Atoi(f[:3])
	if err != nil {
		return nil, err
	}
	if l+3 > len(buf) {
		return nil, errors.New("Wrong size")
	}
	fmt.Println(l)
	msg := buf[3 : l+3]
	data := make([]byte, l)
	_, err = hex.Decode(data, msg)
	if err != nil {
		return nil, err
	}
	return data[:l/2], nil
}
