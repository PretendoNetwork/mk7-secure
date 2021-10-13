package main

import (
	"strconv"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getSessionURLs(err error, client *nex.Client, callID uint32, gatheringId uint32) {
	var stationUrlStrings []string
	for _, gatheringClient := range MatchmakingState[gatheringId].clients {
		if(gatheringClient.PID() == MatchmakingState[gatheringId].matchmakeSession.Gathering.HostPID){
			var localStation nex.StationURL
			localStation.FromString(gatheringClient.LocalStationUrl())
			rvcId := strconv.Itoa(int(gatheringClient.ConnectionId()))
			pid := strconv.Itoa(int(MatchmakingState[gatheringId].matchmakeSession.Gathering.HostPID))
			localStation.SetRVCID(&rvcId)
			localStation.SetPid(&pid)

			stationUrlStrings = append(stationUrlStrings, localStation.EncodeToString())
			var globalStation nex.StationURL
			globalStation.FromString("prudp:/address=0.0.0.0;port=0;PID=0;RVCID=6723910;natf=2;natm=1;pmp=0;sid=15;type=3;upnp=0")

			address := gatheringClient.Address().IP.String()
			port := strconv.Itoa(gatheringClient.Address().Port)

			globalStation.SetPid(&pid)
			globalStation.SetAddress(&address)
			globalStation.SetPort(&port)
			globalStation.SetRVCID(&rvcId)

			stationURL := globalStation.EncodeToString()
			stationUrlStrings = append(stationUrlStrings, stationURL)
		}
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteListString(stationUrlStrings)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingMethodGetSessionURLs, rmcResponseBody)

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
