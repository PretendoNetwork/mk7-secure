package main

import (
	"fmt"
	"strconv"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func requestProbeInitiationExt(err error, client *nex.Client, callID uint32, targetList []string, stationToProbe string) {
	rmcResponse := nex.NewRMCResponse(nexproto.NatTraversalProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.NatTraversalRequestProbeInitiationExt, nil)

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
	
	rmcMessage := nex.RMCRequest{}
	rmcMessage.SetProtocolID(nexproto.NatTraversalProtocolID)
	rmcMessage.SetCallID(0xffff0000+callID)
	rmcMessage.SetMethodID(nexproto.NatTraversalInitiateProbe)
	rmcRequestStream := nex.NewStreamOut(nexServer)
	rmcRequestStream.WriteString(stationToProbe)
	rmcRequestBody := rmcRequestStream.Bytes()
	rmcMessage.SetParameters(rmcRequestBody)
	rmcMessageBytes := rmcMessage.Bytes()
	
	//target := stationToProbe
	for _, target := range targetList {
		var targetStation nex.StationURL
		var targetClient *nex.Client
		var clientFound bool
		targetStation.FromString(target)
		fmt.Println(client.Address().String())
		fmt.Println(targetStation.Address() +":"+ targetStation.Port())
		rvcId, _ := strconv.Atoi(targetStation.RVCID())
		for _, targetClient = range MatchmakingState[0].clients { //need to not hardcode this
			if(targetClient.ConnectionId() == uint32(rvcId)){
				clientFound = true
				break
			}
		}
		if(clientFound){
			//fmt.Println(targetClient.Address().String())

			messagePacket, _ := nex.NewPacketV0(targetClient, nil)
			messagePacket.SetVersion(0)
			messagePacket.SetSource(0xA1)
			messagePacket.SetDestination(0xAF)
			messagePacket.SetType(nex.DataPacket)
			messagePacket.SetPayload(rmcMessageBytes)
		
			messagePacket.AddFlag(nex.FlagNeedsAck)
			messagePacket.AddFlag(nex.FlagReliable)
		
			nexServer.Send(messagePacket)
		}else{
			fmt.Println("not found"+targetStation.Address() +":"+ targetStation.Port())
		}
	}
	//rmcMessage.SetSuccess(nexproto.MessageDeliveryMethodDeliverMessage, oUserMessage)
}
