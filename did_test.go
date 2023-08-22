package core

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDID(t *testing.T) {

	// did
	didStr := "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ",
		did.ID.String())
	require.Equal(t, Mumbai, did.NetworkID)
	require.Equal(t, Polygon, did.Blockchain)

	// readonly did
	didStr = "did:iden3:tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa"

	did, err = ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, "tN4jDinQUdMuJJo6GbVeKPNTPCJ7txyXTWU4T2tJa",
		did.ID.String())
	require.Equal(t, NetworkID(""), did.NetworkID)
	require.Equal(t, Blockchain(""), did.Blockchain)

	require.Equal(t, [2]byte{DIDMethodByte[DIDMethodIden3], 0b0}, did.ID.Type())
}

func TestCreateNewShibDID(t *testing.T) {

	p, _ := BuildDIDType(DIDMethodShib, Shibarium, PuppyNet)
	require.Equal(t, p[0], uint8(0x03))
	require.Equal(t, p[1], uint8(0x42))

	typ0, err := BuildDIDType(DIDMethodShib, Shibarium, PuppyNet)
	require.NoError(t, err)

	genesisState, _ := new(big.Int).SetString("0xC46A8845844a1ec44d18C48a7e3a6f919893C728", 0)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	didStr := did.String()
	require.Equal(t, didStr, "did:shib:shibarium:puppynet:3tCVcahkPK58AGXwZuSgFsM1dmCGezVuwsciBFVYiZ")
	did2, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, did.Method, did2.Method)
	require.Equal(t, did.ID, did2.ID)
}

func TestCreateNewShibReadonlyDID(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodShib, NoChain, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	didStr := did.String()

	require.Equal(t, didStr, "did:shib:3etR8HpjUyg29h5ssFaTghiBezcpmhidcZ6Z5WuNr3")

	did2, err := ParseDID(didStr)
	require.NoError(t, err)

	require.Equal(t, did.Method, did2.Method)
	require.Equal(t, did.ID, did2.ID)
}

func TestParseShibDID(t *testing.T) {

	didStr := "did:shib:shibarium:main:3suph5aVnT3uDycoQqtYjm5eJx2295k3L5XeHXd8Kq"

	did, err := ParseDID(didStr)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, "3suph5aVnT3uDycoQqtYjm5eJx2295k3L5XeHXd8Kq", did.ID.String())
	method := did.Method
	require.Equal(t, DIDMethodShib, method)
	blockchain := did.Blockchain
	require.Equal(t, Shibarium, blockchain)
	networkID := did.NetworkID
	require.Equal(t, Main, networkID)

}

func TestDID_MarshalJSON(t *testing.T) {
	id, err := IDFromString("wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ")
	require.NoError(t, err)
	did := DID{
		ID:         id,
		Method:     DIDMethodIden3,
		Blockchain: Polygon,
		NetworkID:  Mumbai,
	}

	b, err := did.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t,
		`"did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"`,
		string(b))
}

func TestDID_UnmarshalJSON(t *testing.T) {
	inBytes := `{"obj": "did:iden3:polygon:mumbai:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	id, err := IDFromString("wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ")
	require.NoError(t, err)
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err = json.Unmarshal([]byte(inBytes), &obj)
	require.NoError(t, err)
	require.NotNil(t, obj.Obj)
	require.Equal(t, id, obj.Obj.ID)
	require.Equal(t, DIDMethodIden3, obj.Obj.Method)
	require.Equal(t, Polygon, obj.Obj.Blockchain)
	require.Equal(t, Mumbai, obj.Obj.NetworkID)
}

func TestDID_UnmarshalJSON_Error(t *testing.T) {
	inBytes := `{"obj": "did:iden3:eth:goerli:wyFiV4w71QgWPn6bYLsZoysFay66gKtVa9kfu6yMZ"}`
	var obj struct {
		Obj *DID `json:"obj"`
	}
	err := json.Unmarshal([]byte(inBytes), &obj)
	require.EqualError(t, err,
		"invalid did format: network method of core identity mumbai differs from given did network specific id goerli")
}

func TestDIDGenesisFromState(t *testing.T) {

	typ0, err := BuildDIDType(DIDMethodIden3, NoChain, NoNetwork)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ0, genesisState)
	require.NoError(t, err)

	require.Equal(t, DIDMethodIden3, did.Method)
	require.Equal(t, NoChain, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)
	require.Equal(t, "did:iden3:tJ93RwaVfE1PEMxd5rpZZuPtLCwbEaDCrNBhAy8HM", did.String())
}

func TestDID_PolygonID_Types(t *testing.T) {

	// Polygon no chain, no network
	did := helperBuildDIDFromType(t, DIDMethodPolygonID, NoChain, NoNetwork)

	require.Equal(t, DIDMethodPolygonID, did.Method)
	require.Equal(t, NoChain, did.Blockchain)
	require.Equal(t, NoNetwork, did.NetworkID)
	require.Equal(t, "did:polygonid:2mbH5rt9zKT1mTivFAie88onmfQtBU9RQhjNPLwFZh", did.String())

	// Polygon | Polygon chain, Main
	did2 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Main)

	require.Equal(t, DIDMethodPolygonID, did2.Method)
	require.Equal(t, Polygon, did2.Blockchain)
	require.Equal(t, Main, did2.NetworkID)
	require.Equal(t, "did:polygonid:polygon:main:2pzr1wiBm3Qhtq137NNPPDFvdk5xwRsjDFnMxpnYHm", did2.String())

	// Polygon | Polygon chain, Mumbai
	did3 := helperBuildDIDFromType(t, DIDMethodPolygonID, Polygon, Mumbai)

	require.Equal(t, DIDMethodPolygonID, did3.Method)
	require.Equal(t, Polygon, did3.Blockchain)
	require.Equal(t, Mumbai, did3.NetworkID)
	require.Equal(t, "did:polygonid:polygon:mumbai:2qCU58EJgrELNZCDkSU23dQHZsBgAFWLNpNezo1g6b", did3.String())

}

func helperBuildDIDFromType(t testing.TB, method DIDMethod, blockchain Blockchain, network NetworkID) *DID {
	t.Helper()

	typ, err := BuildDIDType(method, blockchain, network)
	require.NoError(t, err)

	genesisState := big.NewInt(1)
	did, err := DIDGenesisFromIdenState(typ, genesisState)
	require.NoError(t, err)

	return did
}
