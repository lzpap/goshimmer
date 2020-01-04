package autopeering

import (
	flag "github.com/spf13/pflag"
)

const (
	CFG_ENTRY_NODES = "autopeering.entryNodes"
	CFG_SELECTION   = "autopeering.selection"
)

func init() {
	flag.StringSlice(CFG_ENTRY_NODES, []string{"V8LYtWWcPYYDTTXLeIEFjJEuWlsjDiI0+Pq/Cx9ai6g=@116.202.49.178:14626"}, "list of trusted entry nodes for auto peering")
	flag.Bool(CFG_SELECTION, true, "enable peer selection")
}
