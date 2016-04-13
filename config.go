package main

var config struct {
	Local, Multicast struct {
		Network, Address string
	}
    Remote []struct {
        Network, Address string
        Cost uint64
    }
	NetworkInterface string
	PrivateKeyPath string
    AdvertiseInterval uint64
    HelloInterval uint64
}