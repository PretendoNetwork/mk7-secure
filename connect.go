package main

import (
	nex "github.com/PretendoNetwork/nex-go"
)

func connect(packet *nex.PacketV0) {
	payload := packet.Payload()

	stream := nex.NewStreamIn(payload, nexServer)

	_, _ = stream.ReadBuffer()
	checkData, _ := stream.ReadBuffer()

	sessionKey := make([]byte, nexServer.KerberosKeySize())

	kerberos := nex.NewKerberosEncryption(sessionKey)

	checkDataDecrypted := kerberos.Decrypt(checkData)
	checkDataStream := nex.NewStreamIn(checkDataDecrypted, nexServer)

	userPID := checkDataStream.ReadUInt32LE() // User PID
	packet.Sender().SetPID(userPID)
	_ = checkDataStream.ReadUInt32LE() //CID of secure server station url
	responseCheck := checkDataStream.ReadUInt32LE()

	responseValueStream := nex.NewStreamOut(nexServer)
	responseValueStream.WriteUInt32LE(responseCheck + 1)

	responseValueBufferStream := nex.NewStreamOut(nexServer)
	responseValueBufferStream.WriteBuffer(responseValueStream.Bytes())

	nexServer.AcknowledgePacket(packet, responseValueBufferStream.Bytes())

	packet.Sender().UpdateRC4Key(sessionKey)
	packet.Sender().SetSessionKey(sessionKey)
}
