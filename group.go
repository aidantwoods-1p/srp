package srp

import (
	"math/big"
)

// Group is a Diffie-Hellman group and has an unexported generator and modulus.
// It has a Label or name that the group can call itself.
// Recommended ExponentSize (in bytes) is based on the
// lower estimates given in section 8 of RFC 3526 for the ephemeral random exponents.
type Group struct {
	g, n         *big.Int
	Label        string
	ExponentSize int // RFC 3526 §8
}

// NewGroup creates and initializes a an SRP group.
func NewGroup() *Group {
	return &Group{
		g: &big.Int{},
		n: &big.Int{},
	}
}

// N returns the modulus of the the group.
func (g *Group) N() *big.Int {
	return g.n
}

// Generator returns little g, the generator for the group as a big int.
func (g *Group) Generator() *big.Int {
	return g.g
}

// RFC 5054 groups are listed by their numbers in Appendix A of the RFC.
const (
	// The values correspond to the numbering in Appendix A of RFC 5054
	// so not using iota mechanism for numbering here.
	RFC5054Group1024 = 1 // We won't allow this group
	RFC5054Group1536 = 2 // We aren't going to allow this one either
	RFC5054Group2048 = 3
	RFC5054Group3072 = 4
	RFC5054Group4096 = 5
	RFC5054Group6144 = 6
	RFC5054Group8192 = 7
)

// KnownGroups is a map from strings to Diffie-Hellman group parameters.
var KnownGroups = make(map[int]*Group)

// MinGroupSize (in bits) sets a lower bound on the size of DH groups
// that will pass certain internal checks. Defaults to 2048.
var MinGroupSize = 2048

// MinExponentSize (in bytes) for generating ephemeral private keys.
var MinExponentSize = 32

func init() {
	g2048n := NumberFromString("0xAC6BDB41324A9A9BF166DE5E1389582FAF72B6651987EE07FC319294" +
		"3DB56050A37329CBB4A099ED8193E0757767A13DD52312AB4B03310D" +
		"CD7F48A9DA04FD50E8083969EDB767B0CF6095179A163AB3661A05FB" +
		"D5FAAAE82918A9962F0B93B855F97993EC975EEAA80D740ADBF4FF74" +
		"7359D041D5C33EA71D281E446B14773BCA97B43A23FB801676BD207A" +
		"436C6481F1D2B9078717461A5B9D32E688F87748544523B524B0D57D" +
		"5EA77A2775D2ECFA032CFBDBF52FB3786160279004E57AE6AF874E73" +
		"03CE53299CCC041C7BC308D82A5698F3A8D0C38271AE35F8E9DBFBB6" +
		"94B5C803D89F7AE435DE236D525F54759B65E372FCD68EF20FA7111F" +
		"9E4AFF73")
	g2048 := &Group{
		g:            big.NewInt(2),
		n:            g2048n,
		Label:        "5054A2048",
		ExponentSize: 27}

	g3072n := NumberFromString("0xFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF")
	g3072 := &Group{
		g:            big.NewInt(5),
		n:            g3072n,
		Label:        "5054A3072",
		ExponentSize: 32}

	// RFC 3526 id 16
	g4096n := NumberFromString("0xFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08" +
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B" +
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9" +
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6" +
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8" +
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C" +
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718" +
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D" +
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D" +
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226" +
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC" +
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26" +
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB" +
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2" +
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127" +
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199" +
		"FFFFFFFFFFFFFFFF")
	g4096 := &Group{
		g:            big.NewInt(5),
		n:            g4096n,
		Label:        "5054A4096",
		ExponentSize: 38}

	// RFC 3526 group id 17
	g6144n := NumberFromString("0xFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08" +
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B" +
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9" +
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6" +
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8" +
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C" +
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718" +
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D" +
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D" +
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226" +
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC" +
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26" +
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB" +
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2" +
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127" +
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934028492" +
		"36C3FAB4D27C7026C1D4DCB2602646DEC9751E763DBA37BDF8FF9406" +
		"AD9E530EE5DB382F413001AEB06A53ED9027D831179727B0865A8918" +
		"DA3EDBEBCF9B14ED44CE6CBACED4BB1BDB7F1447E6CC254B33205151" +
		"2BD7AF426FB8F401378CD2BF5983CA01C64B92ECF032EA15D1721D03" +
		"F482D7CE6E74FEF6D55E702F46980C82B5A84031900B1C9E59E7C97F" +
		"BEC7E8F323A97A7E36CC88BE0F1D45B7FF585AC54BD407B22B4154AA" +
		"CC8F6D7EBF48E1D814CC5ED20F8037E0A79715EEF29BE32806A1D58B" +
		"B7C5DA76F550AA3D8A1FBFF0EB19CCB1A313D55CDA56C9EC2EF29632" +
		"387FE8D76E3C0468043E8F663F4860EE12BF2D5B0B7474D6E694F91E" +
		"6DCC4024FFFFFFFFFFFFFFFF")
	g6144 := &Group{
		g:            big.NewInt(5),
		n:            g6144n,
		Label:        "5054A6144",
		ExponentSize: 43}

	// RFC 3526 group id 18
	g8192n := NumberFromString("0xFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08" +
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B" +
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9" +
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6" +
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8" +
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C" +
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718" +
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D" +
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D" +
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226" +
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC" +
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26" +
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB" +
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2" +
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127" +
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934028492" +
		"36C3FAB4D27C7026C1D4DCB2602646DEC9751E763DBA37BDF8FF9406" +
		"AD9E530EE5DB382F413001AEB06A53ED9027D831179727B0865A8918" +
		"DA3EDBEBCF9B14ED44CE6CBACED4BB1BDB7F1447E6CC254B33205151" +
		"2BD7AF426FB8F401378CD2BF5983CA01C64B92ECF032EA15D1721D03" +
		"F482D7CE6E74FEF6D55E702F46980C82B5A84031900B1C9E59E7C97F" +
		"BEC7E8F323A97A7E36CC88BE0F1D45B7FF585AC54BD407B22B4154AA" +
		"CC8F6D7EBF48E1D814CC5ED20F8037E0A79715EEF29BE32806A1D58B" +
		"B7C5DA76F550AA3D8A1FBFF0EB19CCB1A313D55CDA56C9EC2EF29632" +
		"387FE8D76E3C0468043E8F663F4860EE12BF2D5B0B7474D6E694F91E" +
		"6DBE115974A3926F12FEE5E438777CB6A932DF8CD8BEC4D073B931BA" +
		"3BC832B68D9DD300741FA7BF8AFC47ED2576F6936BA424663AAB639C" +
		"5AE4F5683423B4742BF1C978238F16CBE39D652DE3FDB8BEFC848AD9" +
		"22222E04A4037C0713EB57A81A23F0C73473FC646CEA306B4BCBC886" +
		"2F8385DDFA9D4B7FA2C087E879683303ED5BDD3A062B3CF5B3A278A6" +
		"6D2A13F83F44F82DDF310EE074AB6A364597E899A0255DC164F31CC5" +
		"0846851DF9AB48195DED7EA1B1D510BD7EE74D73FAF36BC31ECFA268" +
		"359046F4EB879F924009438B481C6CD7889A002ED5EE382BC9190DA6" +
		"FC026E479558E4475677E9AA9E3050E2765694DFC81F56E880B96E71" +
		"60C980DD98EDD3DFFFFFFFFFFFFFFFFF")
	g8192 := &Group{
		g:            big.NewInt(19),
		n:            g8192n,
		Label:        "5054A8192",
		ExponentSize: 48}

	KnownGroups[RFC5054Group2048] = g2048
	KnownGroups[RFC5054Group3072] = g3072
	KnownGroups[RFC5054Group4096] = g4096
	KnownGroups[RFC5054Group6144] = g6144
	KnownGroups[RFC5054Group8192] = g8192
	// DefaultGroup := g4096
}

/**
 ** Copyright 2017 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
