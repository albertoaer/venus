package protocol

import (
	"net"
)

type peer struct {
	address       net.Addr
	channel       PacketChannel
	lastTimestamp int64
}

func newPeer(channel PacketChannel, address net.Addr) *peer {
	return &peer{
		channel:       channel,
		address:       address,
		lastTimestamp: 0,
	}
}

func (peer *peer) createPacket(data []byte) Packet {
	return Packet{
		Data:    data,
		Address: peer.address,
		Channel: peer.channel,
	}
}

func (peer *peer) updateLastTimestamp(timestamp int64) bool {
	if timestamp <= peer.lastTimestamp {
		return false
	}
	peer.lastTimestamp = timestamp
	return true
}

type peerManager struct {
	// peers associated by client, is uncommon but can be multiple associated to one client id
	//
	// Iterating this map might produce repeated peers since peers can be aliased multiple times
	clientPeers map[ClientId][]*peer
	// peers associated by address, used to prevent multiple peers from holding the same channel-address
	// with different connection state data like lastTimestamp...
	//
	//	Iterating this map does not produce repeated peers, best option for broadcast
	addressPeers map[string][]*peer
}

func newPeerManager() *peerManager {
	return &peerManager{
		clientPeers:  make(map[ClientId][]*peer),
		addressPeers: make(map[string][]*peer),
	}
}

func (pm *peerManager) addPeer(id ClientId, address net.Addr, channel PacketChannel) *peer {
	strAddress := address.String()
	var targetPeer *peer
	if peers, exists := pm.addressPeers[strAddress]; exists {
		for _, currentPeer := range peers {
			if currentPeer.channel == channel {
				targetPeer = currentPeer
				break
			}
		}
		if targetPeer == nil {
			targetPeer = newPeer(channel, address)
			pm.addressPeers[strAddress] = append(peers, targetPeer)
		}
	} else {
		targetPeer = newPeer(channel, address)
		pm.addressPeers[strAddress] = []*peer{targetPeer}
	}
	if peers, exists := pm.clientPeers[id]; exists {
		found := false
		for _, currentPeer := range peers {
			if currentPeer == targetPeer {
				found = true
			}
		}
		if found {
			pm.clientPeers[id] = append(peers, targetPeer)
		}
	} else {
		pm.clientPeers[id] = []*peer{targetPeer}
	}
	return targetPeer
}

func getPeerContainsComparator(arr []*peer) func(*peer) bool {
	if arr == nil {
		return func(*peer) bool { return false }
	}
	return func(item *peer) bool {
		for _, peer := range arr {
			if peer == item {
				return true
			}
		}
		return false
	}
}

func (pm *peerManager) handlePeerData(
	id ClientId,
	address net.Addr,
	channel PacketChannel,
	timestamp int64,
) (valid bool) {
	var targetPeer *peer
	if peers, exists := pm.clientPeers[id]; exists {
		for _, peer := range peers {
			if peer.address.String() == address.String() && peer.channel == channel {
				targetPeer = peer
			}
		}
	}
	if targetPeer == nil {
		targetPeer = pm.addPeer(id, address, channel)
	}
	valid = targetPeer.updateLastTimestamp(timestamp)
	return
}

func (pm *peerManager) spreadDataByClient(
	receiver *ClientId,
	sender ClientId,
	data []byte,
	broadcastIfNotFound bool,
) (packets []Packet) {
	packets = make([]Packet, 0, 1)
	if receiver != nil {
		if peers, found := pm.clientPeers[*receiver]; found {
			for _, targetPeer := range peers {
				packets = append(packets, targetPeer.createPacket(data))
			}
			return
		}
	}
	if broadcastIfNotFound {
		contains := getPeerContainsComparator(pm.clientPeers[sender])
		for _, peers := range pm.addressPeers {
			for _, targetPeer := range peers {
				if !contains(targetPeer) {
					packets = append(packets, targetPeer.createPacket(data))
				}
			}
		}
	}
	return
}
