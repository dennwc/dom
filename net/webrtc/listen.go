package webrtc

type OfferListener interface {
	Listener
	Answer(data []byte) error
}

func (s *Local) Listen(lis OfferListener) (*Peers, error) {
	c := newPeerConnection()
	return &Peers{l: s, c: c, lis: lis, ans: lis}, nil
}
