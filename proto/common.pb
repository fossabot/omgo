message PbCommon {

// commands
// S(erver) -> C(lient) : even
// C -> S : odd
enum Cmd {
  kHeartBeatReq = 1;
  kHeartBeatRsp = 2;
  kLoginReq = 11;
  kLoginRsp = 12;
  kGetSeedReq = 31;
  kGetSeedRsp = 32;
  kProtoPingReq = 1001;
  kProtoPingRsp = 1002;
}

message RspHeader {
   int32 status = 1;
   string msg = 2;
   uint64 seq = 3;
}
}
