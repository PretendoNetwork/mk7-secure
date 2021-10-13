package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

type MatchmakingData struct {
	matchmakeSession    *nexproto.MatchmakeSession
	clients             []*nex.Client
}

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol
var MatchmakingState []*MatchmakingData

func main() {

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(0)
	nexServer.SetNexVersion(2)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetSignatureVersion(0)
	nexServer.SetAccessKey("6181dff1")

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==MK7 - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("=================")
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	natTraversalProtocolServer := nexproto.NewNatTraversalProtocol(nexServer)
	utilityProtocolServer := nexproto.NewUtilityProtocol(nexServer)
	matchmakeExtensionProtocolServer := nexproto.NewMatchmakeExtensionProtocol(nexServer)
	matchMakingProtocolServer := nexproto.NewMatchMakingProtocol(nexServer)

	//needed for the datastore method MK7 contacts when first going online (just needs a response of some kind)
	dataStorePrococolServer := nexproto.NewDataStoreProtocol(nexServer)
	_ = dataStorePrococolServer

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	secureServer.Register(register)

	natTraversalProtocolServer.RequestProbeInitiationExt(requestProbeInitiationExt)
	natTraversalProtocolServer.ReportNatProperties(reportNatProperties)

	utilityProtocolServer.GetAssociatedNexUniqueIdWithMyPrincipalId(getAssociatedNexUniqueIdWithMyPrincipalId)

	matchmakeExtensionProtocolServer.AutoMatchmake_Postpone(autoMatchmake_Postpone)

	matchMakingProtocolServer.GetSessionURLs(getSessionURLs)

	nexServer.Listen(":60003")
}
