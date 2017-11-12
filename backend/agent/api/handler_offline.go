package api

// ProcOfflineReq handles client offline request
func ProcOfflineReq(session *Session, inPacket []byte) []byte {
	session.SetFlagKicked()
	return nil
}
