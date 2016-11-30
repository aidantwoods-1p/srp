package srp

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

type srpgroup struct {
	g, N *big.Int
}

var knownGroups map[string]*srpgroup

func init() {
	// g1024 is used for unit tests only
	g1024 := &srpgroup{g: big.NewInt(2), N: new(big.Int)}
	g1024.N.SetString("EEAF0AB9ADB38DD69C33F80AFA8FC5E86072618775FF3C0B9EA2314C"+
		"9C256576D674DF7496EA81D3383B4813D692C6E0E0D5D8E250B98BE4"+
		"8E495C1D6089DAD15DC7D7B46154D6B6CE8EF4AD69B15D4982559B29"+
		"7BCF1885C529F566660E57EC68EDBC3C05726CC02FD4CBF4976EAA9A"+
		"FD5138FE8376435B9FC61D2FC0EB06E3", 16)

	g4096 := &srpgroup{g: big.NewInt(5), N: new(big.Int)}
	g4096.N.SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08"+
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B"+
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9"+
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6"+
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8"+
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D"+
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C"+
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718"+
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D"+
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D"+
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226"+
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C"+
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC"+
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26"+
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB"+
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2"+
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127"+
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199"+
		"FFFFFFFFFFFFFFFF", 16)

	g6144 := &srpgroup{g: big.NewInt(5), N: new(big.Int)}
	g6144.N.SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08"+
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B"+
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9"+
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6"+
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8"+
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D"+
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C"+
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718"+
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D"+
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D"+
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226"+
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C"+
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC"+
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26"+
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB"+
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2"+
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127"+
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934028492"+
		"36C3FAB4D27C7026C1D4DCB2602646DEC9751E763DBA37BDF8FF9406"+
		"AD9E530EE5DB382F413001AEB06A53ED9027D831179727B0865A8918"+
		"DA3EDBEBCF9B14ED44CE6CBACED4BB1BDB7F1447E6CC254B33205151"+
		"2BD7AF426FB8F401378CD2BF5983CA01C64B92ECF032EA15D1721D03"+
		"F482D7CE6E74FEF6D55E702F46980C82B5A84031900B1C9E59E7C97F"+
		"BEC7E8F323A97A7E36CC88BE0F1D45B7FF585AC54BD407B22B4154AA"+
		"CC8F6D7EBF48E1D814CC5ED20F8037E0A79715EEF29BE32806A1D58B"+
		"B7C5DA76F550AA3D8A1FBFF0EB19CCB1A313D55CDA56C9EC2EF29632"+
		"387FE8D76E3C0468043E8F663F4860EE12BF2D5B0B7474D6E694F91E"+
		"6DCC4024FFFFFFFFFFFFFFFF", 16)

	g8192 := &srpgroup{g: big.NewInt(19), N: new(big.Int)}
	g8192.N.SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E08"+
		"8A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B"+
		"302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9"+
		"A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE6"+
		"49286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8"+
		"FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D"+
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C"+
		"180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF695581718"+
		"3995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D"+
		"04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7D"+
		"B3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D226"+
		"1AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200C"+
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFC"+
		"E0FD108E4B82D120A92108011A723C12A787E6D788719A10BDBA5B26"+
		"99C327186AF4E23C1A946834B6150BDA2583E9CA2AD44CE8DBBBC2DB"+
		"04DE8EF92E8EFC141FBECAA6287C59474E6BC05D99B2964FA090C3A2"+
		"233BA186515BE7ED1F612970CEE2D7AFB81BDD762170481CD0069127"+
		"D5B05AA993B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934028492"+
		"36C3FAB4D27C7026C1D4DCB2602646DEC9751E763DBA37BDF8FF9406"+
		"AD9E530EE5DB382F413001AEB06A53ED9027D831179727B0865A8918"+
		"DA3EDBEBCF9B14ED44CE6CBACED4BB1BDB7F1447E6CC254B33205151"+
		"2BD7AF426FB8F401378CD2BF5983CA01C64B92ECF032EA15D1721D03"+
		"F482D7CE6E74FEF6D55E702F46980C82B5A84031900B1C9E59E7C97F"+
		"BEC7E8F323A97A7E36CC88BE0F1D45B7FF585AC54BD407B22B4154AA"+
		"CC8F6D7EBF48E1D814CC5ED20F8037E0A79715EEF29BE32806A1D58B"+
		"B7C5DA76F550AA3D8A1FBFF0EB19CCB1A313D55CDA56C9EC2EF29632"+
		"387FE8D76E3C0468043E8F663F4860EE12BF2D5B0B7474D6E694F91E"+
		"6DBE115974A3926F12FEE5E438777CB6A932DF8CD8BEC4D073B931BA"+
		"3BC832B68D9DD300741FA7BF8AFC47ED2576F6936BA424663AAB639C"+
		"5AE4F5683423B4742BF1C978238F16CBE39D652DE3FDB8BEFC848AD9"+
		"22222E04A4037C0713EB57A81A23F0C73473FC646CEA306B4BCBC886"+
		"2F8385DDFA9D4B7FA2C087E879683303ED5BDD3A062B3CF5B3A278A6"+
		"6D2A13F83F44F82DDF310EE074AB6A364597E899A0255DC164F31CC5"+
		"0846851DF9AB48195DED7EA1B1D510BD7EE74D73FAF36BC31ECFA268"+
		"359046F4EB879F924009438B481C6CD7889A002ED5EE382BC9190DA6"+
		"FC026E479558E4475677E9AA9E3050E2765694DFC81F56E880B96E71"+
		"60C980DD98EDD3DFFFFFFFFFFFFFFFFF", 16)

	knownGroups = make(map[string]*srpgroup)
	knownGroups["1024"] = g1024
	knownGroups["4096"] = g4096
	knownGroups["6144"] = g6144
	knownGroups["8192"] = g8192
}

// AmodNisValid determines if "A mod N" is valid for the given
// SRP group and value of A.
func AmodNisValid(A *big.Int, groupName string) bool {
	result := big.Int{}

	group := knownGroups[groupName]
	if group == nil {
		return false
	}

	result.Mod(A, group.N)
	if result.Sign() == 0 { // sign is zero only when the whole value is 0.
		return false
	}
	return true
}

// CalculateVerifier calculates the verifier
func CalculateVerifier(groupName string, x *big.Int) *big.Int {
	group := knownGroups[groupName]

	i := new(big.Int)
	return i.Exp(group.g, x, group.N)
}

// prehash is kept for compatibility with legacy implementations
func prehash(s string) string {
	if s == "" {
		return ""
	}

	hasher := sha256.New()
	hasher.Write([]byte(s))
	bits := hasher.Sum(nil)

	return strings.TrimRight(base32.StdEncoding.EncodeToString(bits), "=")
}

// bytesToHex returns hexadecimal representation of the slice.
func bytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// CalculateX compute X value used in SRP authentication.
func CalculateX(method, alg, email, password string, salt []byte, iterations int, accountKey *AccountKey) (*big.Int, error) {
	if iterations == 0 { // Using SRP Test Vectors
		h1 := sha1.New()
		h1.Write(salt)

		h2 := sha1.New()
		h2.Write([]byte(email + ":" + password))
		h1.Write(h2.Sum(nil))

		return NumberFromBytes(h1.Sum(nil)), nil
	}

	if accountKey == nil {
		return nil, errors.New("missing AccountKey in CalculateX")
	}

	var h func() hash.Hash
	var keyLen int
	var err error
	salt, err = base64.RawURLEncoding.DecodeString(string(salt))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode salt")
	}

	if alg == "PBES2-HS512" || alg == "PBES2g-HS512" {
		keyLen = 512 / 8
		h = sha512.New
	} else if alg == "PBES2-HS256" || alg == "PBES2g-HS256" {
		keyLen = 256 / 8
		h = sha256.New
	} else {
		return nil, fmt.Errorf("invalid SRP alg: %q", alg)
	}

	if strings.HasPrefix(method, "SRP-") {
		derivedBits := pbkdf2.Key([]byte(prehash(password)), salt, iterations, keyLen, h)
		combined := accountKey.CombineWithBytes(derivedBits)

		hasher := sha1.New()

		hasher.Write(salt)
		hasher.Write([]byte(email + ":" + bytesToHex(combined)))
		return NumberFromBytes(hasher.Sum(nil)), nil
	}

	if strings.HasPrefix(method, "SRPg-") {
		emailSalt := []byte(email)
		info := []byte(method)
		bigSalt := make([]byte, 32)
		if _, err := io.ReadFull(hkdf.New(sha256.New, salt, emailSalt, info), bigSalt); err != nil {
			return nil, errors.Wrap(err, "HKDF failed")
		}

		derivedBits := pbkdf2.Key([]byte(password), bigSalt, iterations, keyLen, h)
		combined := accountKey.CombineWithBytes(derivedBits)

		return NumberFromBytes(combined), nil
	}

	return nil, fmt.Errorf("invalid SRP method: %q", method)
}

// CalculateA computes SRP A value based on a. The a should be randomly generated using `srp.RandomNumber(32)`
func CalculateA(groupName string, a *big.Int) *big.Int {
	group := knownGroups[groupName]
	result := new(big.Int)
	return result.Exp(group.g, a, group.N)
}

// CalculateB calculates B according to SRP RFC
func CalculateB(groupName string, k *big.Int, v *big.Int, randomKey *big.Int) *big.Int {
	group := knownGroups[groupName]

	result := new(big.Int)
	result.Exp(group.g, randomKey, group.N)

	m := new(big.Int)
	m.Mul(k, v)

	result.Add(m, result)
	return result.Mod(result, group.N)
}

// CalculateClientRawKey calculates the raw key
func CalculateClientRawKey(groupName string, a, b, u, x, k *big.Int) *big.Int {
	group := knownGroups[groupName]

	p := new(big.Int)
	r := new(big.Int)
	r.Mul(u, x)
	p.Add(a, r)
	base := new(big.Int)
	r1 := new(big.Int)
	r1.Exp(group.g, x, group.N)
	r = new(big.Int)
	r.Mul(r1, k)
	base.Sub(b, r)
	result := new(big.Int)
	result.Exp(base, p, group.N)

	hex := fmt.Sprintf("%x", result)

	hasher := sha256.New()
	hasher.Write([]byte(hex))
	return NumberFromBytes(hasher.Sum(nil))
}

// CalculateRawKey calculates the raw key
func CalculateRawKey(groupName string, A, v, b, u *big.Int) *big.Int {
	group := knownGroups[groupName]

	result := new(big.Int)
	result.Exp(v, u, group.N)
	result.Mul(result, A)
	return result.Exp(result, b, group.N)
}

// NumberFromString converts a string to a number
func NumberFromString(s string) *big.Int {
	n := strings.Replace(s, " ", "", -1)

	result := new(big.Int)
	result.SetString(strings.TrimPrefix(n, "0x"), 16)

	return result
}

// NumberFromBytes converts a byte array to a number
func NumberFromBytes(bytes []byte) *big.Int {
	result := new(big.Int)
	for _, b := range bytes {
		result.Lsh(result, 8)
		result.Add(result, big.NewInt(int64(b)))
	}

	return result
}

// RandomNumber returns a random number
func RandomNumber() *big.Int {
	bytes := make([]byte, 8)
	rand.Read(bytes)

	return NumberFromBytes(bytes)
}
