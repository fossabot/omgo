package api

// ProcOfflineReq handles client offline request
func ProcOfflineReq(session *Session, inPacket *IncomingPacket) []byte {
	session.SetFlagKicked()
	return nil
}
