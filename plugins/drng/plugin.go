package drng

import (
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/node"
)

// region plugin module setup //////////////////////////////////////////////////////////////////////////////////////////

var PLUGIN = node.NewPlugin("DRNG", node.Enabled, configure, run)
var log *logger.Logger

func configure(*node.Plugin) {
	log = logger.NewLogger("DRNG")
	configureDRNG()
}

func run(*node.Plugin) {
}
