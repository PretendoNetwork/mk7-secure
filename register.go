package main

import (
	"strconv"
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func register(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL) {
	localStation := stationUrls[0]
	localStationURL := localStation.EncodeToString()
	connectionId := uint32(secureServer.ConnectionIDCounter.Increment())
	client.SetConnectionId(connectionId)
	client.SetLocalStationUrl(localStationURL)

	address := client.Address().IP.String()
	port := strconv.Itoa(client.Address().Port)
	natf := "2"
	natm := "1"
	type_ := "3"

	localStation.SetAddress(&address)
	localStation.SetPort(&port)
	localStation.SetNatf(&natf)
	localStation.SetNatm(&natm)
	localStation.SetType(&type_)

	globalStationURL := localStation.EncodeToString()

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt32LE(0x10001) // Success
	rmcResponseStream.WriteUInt32LE(connectionId)
	rmcResponseStream.WriteString(globalStationURL)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.SecureProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureMethodRegister, rmcResponseBody)

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
