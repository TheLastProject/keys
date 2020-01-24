package keys

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatement(t *testing.T) {
	clock := newClock()
	sk := NewEd25519KeyFromSeed(Bytes32(bytes.Repeat([]byte{0x01}, 32)))

	sc := NewSigchain(sk.PublicKey())
	require.Equal(t, 0, sc.Length())
	st, err := GenerateStatement(sc, bytes.Repeat([]byte{0x01}, 16), sk, "test", clock.Now())
	require.NoError(t, err)

	st2, err := NewStatement(st.Sig, st.Data, sk.PublicKey(), st.Seq, st.Prev, st.Revoke, st.Type, st.Timestamp)
	require.NoError(t, err)
	require.Equal(t, st.Bytes(), st2.Bytes())

	rk := GenerateEd25519Key()
	_, err = NewStatement(st.Sig, st.Data, rk.PublicKey(), st.Seq, st.Prev, st.Revoke, st.Type, st.Timestamp)
	require.EqualError(t, err, "verify failed")
}

func TestStatementJSON(t *testing.T) {
	clock := newClock()
	sk := NewEd25519KeyFromSeed(Bytes32(bytes.Repeat([]byte{0x01}, 32)))

	sc := NewSigchain(sk.PublicKey())
	require.Equal(t, 0, sc.Length())

	st, err := GenerateStatement(sc, bytes.Repeat([]byte{0x01}, 16), sk, "test", clock.Now())
	require.NoError(t, err)

	b, err := json.Marshal(st)
	require.NoError(t, err)

	var stOut Statement
	err = json.Unmarshal(b, &stOut)
	require.NoError(t, err)

	require.Equal(t, st.Data, stOut.Data)
	require.Equal(t, st.KID, stOut.KID)
	require.Equal(t, st.Seq, stOut.Seq)
	require.Equal(t, st.Prev, stOut.Prev)
	require.Equal(t, st.Revoke, stOut.Revoke)
	require.Equal(t, st.Type, stOut.Type)
	require.Equal(t, st.serialized, stOut.serialized)

	err = sc.Add(st)
	require.NoError(t, err)

	// Revoke
	revoke, err := GenerateRevoke(sc, 1, sk)
	require.NoError(t, err)

	b2, err := json.Marshal(revoke)
	require.NoError(t, err)

	var stOut2 Statement
	err = json.Unmarshal(b2, &stOut2)
	require.NoError(t, err)

	require.Equal(t, revoke.Data, stOut2.Data)
	require.Equal(t, revoke.KID, stOut2.KID)
	require.Equal(t, revoke.Seq, stOut2.Seq)
	require.Equal(t, revoke.Prev, stOut2.Prev)
	require.Equal(t, revoke.Revoke, stOut2.Revoke)
	require.Equal(t, st.Timestamp, stOut.Timestamp)
	require.Equal(t, revoke.Type, stOut2.Type)
	require.Equal(t, revoke.serialized, stOut2.serialized)
}

func TestStatementSpecificSerialization(t *testing.T) {
	clock := newClock()
	sk := NewEd25519KeyFromSeed(Bytes32(bytes.Repeat([]byte{0x01}, 32)))
	sc := NewSigchain(sk.PublicKey())
	require.Equal(t, 0, sc.Length())

	st, err := GenerateStatement(sc, bytes.Repeat([]byte{0x01}, 16), sk, "test", clock.Now())
	require.NoError(t, err)

	data := statementBytesToSign(st)
	expected := `{".sig":"","data":"AQEBAQEBAQEBAQEBAQEBAQ==","kid":"kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen","seq":1,"ts":1234567890001,"type":"test"}`
	require.Equal(t, expected, string(data))

	dataOut := st.Bytes()
	expectedOut := `{".sig":"/KL3VWwwPnz0rJjWliMFiACgRfMYeHlgwvf0zMtWNo/WXIsuDFNy/PA4LrZqEIWGlkJASfLrGue2kjtCgdjmAA==","data":"AQEBAQEBAQEBAQEBAQEBAQ==","kid":"kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen","seq":1,"ts":1234567890001,"type":"test"}`
	require.Equal(t, expectedOut, string(dataOut))

	require.Equal(t, expectedOut, string(st.Bytes()))

	stOut, err := StatementFromBytes(dataOut)
	require.NoError(t, err)
	require.Equal(t, st.Data, stOut.Data)
	require.Equal(t, st.KID, stOut.KID)
	require.Equal(t, st.Seq, stOut.Seq)
	require.Equal(t, st.Prev, stOut.Prev)
	require.Equal(t, st.Revoke, stOut.Revoke)
	require.Equal(t, st.Type, stOut.Type)
	require.Equal(t, st.serialized, stOut.serialized)

	_, err = StatementFromBytes([]byte("{}"))
	require.EqualError(t, err, "not enough bytes for statement")

	_, err = StatementFromBytes(data)
	require.EqualError(t, err, "statement bytes don't match specific serialization")

	err = sc.Add(st)
	require.NoError(t, err)

	// Revoke
	revoke, err := GenerateRevoke(sc, 1, sk)
	require.NoError(t, err)

	data2 := statementBytesToSign(revoke)
	expected2 := `{".sig":"","kid":"kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen","prev":"tVGumGOvb+QzNfIBuLl6BnKzIOP7PPOaDu+bRrABc58=","revoke":1,"seq":2,"type":"revoke"}`
	require.Equal(t, expected2, string(data2))

	dataOut2 := revoke.Bytes()
	expectedOut2 := `{".sig":"Gh8TJKJM4JHLvaDE/ZhHrKcWVl2sYMFm28Wu9a/ZyEkm2xcaPLKozGCzL8qX97X2b7juVY8ix994Dr4arcJUCQ==","kid":"kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen","prev":"tVGumGOvb+QzNfIBuLl6BnKzIOP7PPOaDu+bRrABc58=","revoke":1,"seq":2,"type":"revoke"}`
	require.Equal(t, expectedOut2, string(dataOut2))

	require.Equal(t, expectedOut2, string(revoke.Bytes()))

	stOut2, stOutErr2 := StatementFromBytes(dataOut2)
	require.NoError(t, stOutErr2)
	require.Equal(t, revoke.Data, stOut2.Data)
	require.Equal(t, revoke.KID, stOut2.KID)
	require.Equal(t, revoke.Seq, stOut2.Seq)
	require.Equal(t, revoke.Prev, stOut2.Prev)
	require.Equal(t, revoke.Revoke, stOut2.Revoke)
	require.Equal(t, revoke.Type, stOut2.Type)
	require.Equal(t, revoke.serialized, stOut2.serialized)
}

func TestStatementKeyURL(t *testing.T) {
	clock := newClock()
	sk := NewEd25519KeyFromSeed(Bytes32(bytes.Repeat([]byte{0x01}, 32)))

	sc := NewSigchain(sk.PublicKey())
	require.Equal(t, 0, sc.Length())

	st, err := GenerateStatement(sc, bytes.Repeat([]byte{0x01}, 16), sk, "test", clock.Now())
	require.NoError(t, err)

	require.Equal(t, "kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen-000000000000001", st.Key())
	require.Equal(t, "/kpe132yw8ht5p8cetl2jmvknewjawt9xwzdlrk2pyxlnwjyqrdq0dawqlrnuen/1", st.URL())
}
