package data

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"log"
	"math"
)

// ForwardEvents is a list of lnrpc.ForwardEvents
type ForwardEvents []*lnrpc.ForwardingEvent

// Node is contains useful meta-data for a remote peer.
type Node struct {
	// PublicKey is the public key of a remote peer.
	PublicKey string
	// Alias is the alias of a remote peer.
	Alias string
}

// Nodes maps chanID to Node.
type Nodes map[uint64]*Node

// BuildData communicates with lnrpc.LightningClient to fetch forwarding events and organize data.
// Returns a list of forwarding events for the node and a map of channel IDs to Node type.
func BuildData(client lnrpc.LightningClient) ([]*lnrpc.ForwardingEvent, Nodes, error) {
	history, err := client.ForwardingHistory(context.Background(), &lnrpc.ForwardingHistoryRequest{
		StartTime: 1514764800,
		//EndTime:              0, // TODO, include endTime
		NumMaxEvents: math.MaxUint32,
	})

	if err != nil {
		return nil, nil, err
	}

	nodes := make(Nodes)

	for _, event := range history.GetForwardingEvents() {

		chanID := event.GetChanIdOut()

		if _, exists := nodes[chanID]; !exists {
			chanResp, err := client.GetChanInfo(context.Background(), &lnrpc.ChanInfoRequest{
				ChanId: chanID,
			})
			if err != nil {
				log.Printf("cannot get channel info for %d\n", chanID)
				continue
			}

			nodePubKey := chanResp.GetNode2Pub()

			nodeResp, err := client.GetNodeInfo(context.Background(), &lnrpc.NodeInfoRequest{
				PubKey:          nodePubKey,
				IncludeChannels: true,
			})

			if err != nil {
				return nil, nil, err
			}

			node := nodeResp.GetNode()

			nodes[chanID] = &Node{
				PublicKey: node.GetPubKey(),
				Alias:     node.GetAlias(),
			}
		}
	}

	return history.GetForwardingEvents(), nodes, nil
}
