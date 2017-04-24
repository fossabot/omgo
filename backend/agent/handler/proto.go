package handler

import "github.com/master-g/omgo/net/packet"

type PkgAutoID struct {
	ID int32
}

func (p PkgAutoID) Pack(w *packet.RawPacket) {
	w.WriteS32(p.ID)
}

type PkgErrorInfo struct {
	Code int32
	Msg  string
}

func (p PkgErrorInfo) Pack(w *packet.RawPacket) {
	w.WriteS32(p.Code)
	w.WriteString(p.Msg)
}

type PkgUserLoginInfo struct {
	MailAddress string
}

func (p PkgUserLoginInfo) Pack(w *packet.RawPacket) {
	w.WriteString(p.MailAddress)
}

type PkgSeedInfo struct {
	ClientSendSeed int32
	ClientRecvSeed int32
}

func (p PkgSeedInfo) Pack(w *packet.RawPacket) {
	w.WriteS32(p.ClientSendSeed)
	w.WriteS32(p.ClientRecvSeed)
}

type PkgUserSnapshot struct {
	UserID int32
}

func (p PkgUserSnapshot) Pack(w *packet.RawPacket) {
	w.WriteS32(p.UserID)
}

func PacketReadAutoID(reader *packet.RawPacket) (tbl PkgAutoID, err error) {
	tbl.ID, err = reader.ReadS32()
	checkErr(err)

	return
}

func PacketReadErrorInfo(reader *packet.RawPacket) (tbl PkgErrorInfo, err error) {
	tbl.Code, err = reader.ReadS32()
	checkErr(err)

	tbl.Msg, err = reader.ReadString()
	checkErr(err)

	return
}

func PacketReadUserLoginInfo(reader *packet.RawPacket) (tbl PkgUserLoginInfo, err error) {
	tbl.MailAddress, err = reader.ReadString()
	checkErr(err)

	return
}

func PacketReadSeedInfo(reader *packet.RawPacket) (tbl PkgSeedInfo, err error) {
	tbl.ClientSendSeed, err = reader.ReadS32()
	checkErr(err)

	tbl.ClientRecvSeed, err = reader.ReadS32()
	checkErr(err)

	return
}

func PacketReadSnapshot(reader *packet.RawPacket) (tbl PkgUserSnapshot, err error) {
	tbl.UserID, err = reader.ReadS32()
	checkErr(err)

	return
}

func checkErr(err error) {
	if err != nil {
		panic("error occured in protocol module")
	}
}
