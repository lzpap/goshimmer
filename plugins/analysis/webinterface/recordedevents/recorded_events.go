package recordedevents

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/iotaledger/goshimmer/plugins/analysis/server"
	"github.com/iotaledger/goshimmer/plugins/analysis/types/heartbeat"
	"github.com/iotaledger/goshimmer/plugins/analysis/webinterface/types"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/node"
)

// Maps nodeId to the latest arrival of a heartbeat
var nodes = make(map[string]time.Time)

// Maps nodeId to outgoing connections + latest arrival of heartbeat
var links = make(map[string]map[string]time.Time)

var lock sync.Mutex

func Configure(plugin *node.Plugin) {
	server.Events.Heartbeat.Attach(events.NewClosure(func(packet heartbeat.Packet) {
		out := ""
		for _, value := range packet.OutboundIDs {
			out += hex.EncodeToString(value)
		}
		in := ""
		for _, value := range packet.InboundIDs {
			in += hex.EncodeToString(value)
		}
		plugin.Node.Logger.Debugw(
			"Heartbeat",
			"nodeId", hex.EncodeToString(packet.OwnID),
			"outboundIds", out,
			"inboundIds", in,
		)
		lock.Lock()
		defer lock.Unlock()

		// process the packet
		nodeIdString := hex.EncodeToString(packet.OwnID)
		timestamp := time.Now()

		// When it is present in the list, we just update the timestamp
		if _, isAlready := nodes[nodeIdString]; !isAlready {
			server.Events.AddNode.Trigger(nodeIdString)
			server.Events.NodeOnline.Trigger(nodeIdString)
		}
		nodes[nodeIdString] = timestamp

		// Outgoing neighbor links update
		for _, outgoingNeighbor := range packet.OutboundIDs {
			outgoingNeighborString := hex.EncodeToString(outgoingNeighbor)
			// Do we already know about this neighbor?
			// If no, add it and set it online
			if _, isAlready := nodes[outgoingNeighborString]; !isAlready {
				// First time we see this particular node
				server.Events.AddNode.Trigger(outgoingNeighborString)
				server.Events.NodeOnline.Trigger(outgoingNeighborString)
			}
			// We have indirectly heard about the neighbor.
			nodes[outgoingNeighborString] = timestamp

			// Update graph when connection hasn't been seen before
			if _, isAlready := links[nodeIdString][outgoingNeighborString]; !isAlready {
				server.Events.ConnectNodes.Trigger(nodeIdString, outgoingNeighborString)
			}
			// Update timestamp
			links[nodeIdString][outgoingNeighborString] = timestamp
		}
		// Incoming neighbor links update
		for _, incomingNeighbor := range packet.OutboundIDs {
			incomingNeighborString := hex.EncodeToString(incomingNeighbor)
			// Do we already know about this neighbor?
			// If no, add it and set it online
			if _, isAlready := nodes[incomingNeighborString]; !isAlready {
				// First time we see this particular node
				server.Events.AddNode.Trigger(incomingNeighborString)
				server.Events.NodeOnline.Trigger(incomingNeighborString)
			}
			// We have indirectly heard about the neighbor.
			nodes[incomingNeighborString] = timestamp

			if _, isAlready := links[incomingNeighborString][nodeIdString]; !isAlready {
				server.Events.ConnectNodes.Trigger(incomingNeighborString, nodeIdString)
			}
			links[incomingNeighborString][nodeIdString] = timestamp
		}
	}))

	go cleanUpPeriodically(CLEAN_UP_PERIOD)
}

// Remove nodes and links we haven't seen for at least 2 times the heartbeat interval
func cleanUpPeriodically(interval time.Duration) {
	for {
		lock.Lock()
		now := time.Now()

		// Go through the list of connections. Remove connections that are older than interval time.
		for srcNode, targetMap := range links {
			for _, lastSeen := range targetMap {
				if now.Sub(lastSeen) > interval {
					delete(links, srcNode)
					server.Events.DisconnectNodes.Trigger()
				}
			}
		}

		// Go through the list of nodes. Remove nodes that haven't been seen for interval time
		for node, lastSeen := range nodes {
			if now.Sub(lastSeen) > interval {
				delete(nodes, node)
				server.Events.NodeOffline.Trigger(node)
				server.Events.RemoveNode.Trigger(node)
			}
		}
		lock.Unlock()
		// Sleep for interval time
		time.Sleep(interval)
	}

}

func getEventsToReplay() (map[string]time.Time, map[string]map[string]time.Time) {
	lock.Lock()
	defer lock.Unlock()

	copiedNodes := make(map[string]time.Time)
	for nodeId, lastHeartbeat := range nodes {
		copiedNodes[nodeId] = lastHeartbeat
	}

	copiedLinks := make(map[string]map[string]time.Time)
	for sourceId, targetMap := range links {
		copiedLinks[sourceId] = make(map[string]time.Time)
		for targetId, lastHeartbeat := range targetMap {
			copiedLinks[sourceId][targetId] = lastHeartbeat
		}
	}

	return copiedNodes, copiedLinks
}

func Replay(handlers *types.EventHandlers) {
	copiedNodes, copiedLinks := getEventsToReplay()

	// When a node is present in the list, it means we heard about it directly
	// or indirectly, but within CLEAN_UP_PERIOD, therefore it is online
	for nodeId, _ := range copiedNodes {
		handlers.AddNode(nodeId)
		handlers.NodeOnline(nodeId)
	}

	for sourceId, targetMap := range copiedLinks {
		for targetId := range targetMap {
			handlers.ConnectNodes(sourceId, targetId)
		}
	}
}
