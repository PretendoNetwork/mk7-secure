package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getAssociatedNexUniqueIdWithMyPrincipalId(err error, client *nex.Client, callID uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt64LE(0)
	rmcResponseStream.WriteUInt64LE(0)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.UtilityProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.UtilityMethodGetAssociatedNexUniqueIdWithMyPrincipalId, rmcResponseBody)

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
