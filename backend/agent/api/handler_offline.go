package api

// ProcOfflineReq handles client offline request
func ProcOfflineReq(session *Session, inPacket *IncomingPacket) *OutgoingPacket {
	session.SetFlagKicked()
	return nil
}
