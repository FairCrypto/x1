package bls

import (
	"encoding/json"
	"errors"
	"math/big"

	bn256 "github.com/umbracle/go-eth-bn256"
)

type PublicKey struct {
	p *bn256.G2
}

// aggregate adds the given public keys
func (p *PublicKey) aggregate(onemore *PublicKey) *PublicKey {
	var g2 *bn256.G2
	if p.p == nil {
		g2 = new(bn256.G2).Set(&zeroG2)
	} else {
		g2 = new(bn256.G2).Set(p.p)
	}

	g2.Add(g2, onemore.p)

	return &PublicKey{p: g2}
}

// Marshal marshal the key to bytes.
func (p *PublicKey) Marshal() []byte {
	return p.p.Marshal()
}

// MarshalJSON implements the json.Marshaler interface.
func (p *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.p.Marshal())
}

// UnmarshalJSON implements the json.Marshaler interface.
func (p *PublicKey) UnmarshalJSON(b []byte) error {
	var d []byte

	err := json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	pk, err := UnmarshalPublicKey(d)
	p.p = pk.p

	return err
}

// UnmarshalPublicKey reads the public key from the given byte array
func UnmarshalPublicKey(raw []byte) (*PublicKey, error) {
	if len(raw) == 0 {
		return nil, errors.New("cannot unmarshal public key from empty slice")
	}

	p := new(bn256.G2)
	_, err := p.Unmarshal(raw)

	return &PublicKey{p: p}, err
}

// ToBigInt converts public key to 4 big ints
func (p PublicKey) ToBigInt() ([4]*big.Int, error) {
	blsKey := p.Marshal()
	res := [4]*big.Int{
		new(big.Int).SetBytes(blsKey[32:64]),
		new(big.Int).SetBytes(blsKey[0:32]),
		new(big.Int).SetBytes(blsKey[96:128]),
		new(big.Int).SetBytes(blsKey[64:96]),
	}

	return res, nil
}

// UnmarshalPublicKeyFromBigInt unmarshals public key from 4 big ints
// Order of coordinates is [A.Y, A.X, B.Y, B.X]
func UnmarshalPublicKeyFromBigInt(b [4]*big.Int) (*PublicKey, error) {
	var pubKeyBuf []byte

	pt1, err := leftPadTo32Bytes(b[1].Bytes())

	if err != nil {
		return nil, err
	}

	pt2, err := leftPadTo32Bytes(b[0].Bytes())

	if err != nil {
		return nil, err
	}

	pt3, err := leftPadTo32Bytes(b[3].Bytes())

	if err != nil {
		return nil, err
	}

	pt4, err := leftPadTo32Bytes(b[2].Bytes())

	if err != nil {
		return nil, err
	}

	pubKeyBuf = append(pubKeyBuf, pt1...)
	pubKeyBuf = append(pubKeyBuf, pt2...)
	pubKeyBuf = append(pubKeyBuf, pt3...)
	pubKeyBuf = append(pubKeyBuf, pt4...)

	return UnmarshalPublicKey(pubKeyBuf)
}

// aggregatePublicKeys calculates P1 + P2 + ...
func aggregatePublicKeys(pubs []*PublicKey) *PublicKey {
	res := *new(bn256.G2).Set(&zeroG2)
	for i := 0; i < len(pubs); i++ {
		res.Add(&res, pubs[i].p)
	}

	return &PublicKey{p: &res}
}