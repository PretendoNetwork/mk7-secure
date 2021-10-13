package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"fmt"
)

func autoMatchmake_Postpone(err error, client *nex.Client, callID uint32, matchmakeSession *nexproto.MatchmakeSession, message string) {
	var foundSession *MatchmakingData
	var foundIndex int
	var element *MatchmakingData
	rmcResponseStream := nex.NewStreamOut(nexServer)
	for foundIndex, element = range MatchmakingState {
		if(uint16(len(element.clients)) < element.matchmakeSession.Gathering.MaximumParticipants ){
			foundSession = element
		}
	}
	if(foundSession == nil){
		foundSession = new(MatchmakingData)
		foundIndex = len(MatchmakingState)
		matchmakeSession.Gathering.ID = uint32(foundIndex)
		matchmakeSession.Gathering.OwnerPID = client.PID()
		matchmakeSession.Gathering.HostPID = client.PID()
		foundSession.matchmakeSession = matchmakeSession
		foundSession.clients = make([]*nex.Client, 0, 0)
	}
	fmt.Println(foundSession.matchmakeSession.Gathering.OwnerPID)
	foundSession.clients = append(foundSession.clients, client)
	MatchmakingState = append(MatchmakingState, foundSession)
	rmcResponseStream.WriteString("MatchmakeSession")
	matchmakeSessionLength := uint32(len(matchmakeSession.Bytes(nex.NewStreamOut(nexServer))))
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength+4)
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength)
	rmcResponseStream.WriteStructure(foundSession.matchmakeSession)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodAutoMatchmake_Postpone, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
