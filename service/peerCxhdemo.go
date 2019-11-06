package service

import (
	"os"
	"strings"

	"github.com/hashicorp/memberlist"

)

type peerCxhdemo struct {
	peer
}

//CreatePeerCxhdemo 为peerCxhdemo工厂
func CreatePeerCxhdemo(
	bindPort int,
	advertisePort int,
	hostname string,
	knownPeers string,
) (*peerCxhdemo, error) {

	p := &peerCxhdemo{}

	kp := os.Getenv("KNOWN_PEERS")
	if len(kp) > 0 {
		p.seedPeers = strings.Split(kp, ";")
	} else {
		p.seedPeers = strings.Split(knownPeers, ";")
	}

	if os.Getenv("PEER_NAME") != "" {
		hostname = os.Getenv("PEER_NAME")
	}

	delegate := newDelegateCxhdemo(p.mlist)

	config := memberlist.DefaultLANConfig()
	config.Delegate = delegate
	config.Name = hostname
	config.BindPort = bindPort
	if advertisePort != 0 {
		config.AdvertisePort = advertisePort
	} else {
		config.AdvertisePort = bindPort
	}
	ml, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	p.mlist = ml

	return p, nil

}

