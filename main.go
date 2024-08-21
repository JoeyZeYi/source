package main

import (
	"github.com/JoeyZeYi/source/excel"
)

type InvitationDataExcel struct {
	Uid           string `excel:"用户ID"`
	NickName      string `excel:"昵称"`
	Sex           string `excel:"性别"`
	RegisterTime  string `excel:"注册时间"`
	ISPhone       string `excel:"手机认证"`
	IsAuth        string `excel:"实名认证"`
	IsHuMan       string `excel:"真人认证"`
	LastLoginTime string `excel:"最后一次登录时间"`
	RechargeRmb   string `excel:"充值金额"` //user_pay_success_log
	//usr_privatemessage_log查这个表
	SayHiNum              string `excel:"主动搭讪次数"`
	SayHiReplyNum         string `excel:"主动搭讪回复次数"`
	SayHiRate             string `excel:"主动搭讪回复率"`
	PassiveSayHiNum       string `excel:"被搭讪次数"`
	PassiveSayHiReplyNum  string `excel:"被搭讪回复次数"`
	PassiveSayHiRate      string `excel:"被搭讪回复率"`
	TxtMsgNum             string `excel:"主动发消息次数"`
	TxtMsgReplyNum        string `excel:"主动发消息回复次数"`
	TxtMsgNumRate         string `excel:"主动发消息回复率"`
	PassiveTxtMsgNum      string `excel:"被动消息次数"`
	PassiveTxtMsgReplyNum string `excel:"被动消息回复次数"`
	PassiveTxtMsgNumRate  string `excel:"被动消息回复率"`
	//chat_room_req_connect_num_log表和chat_room_pri_live_log表
	VoiceNum              string `excel:"主动语音次数"`
	VoiceRefuseNum        string `excel:"主动语音被拒接次数"`
	VoiceRate             string `excel:"主动语音接听率"`
	VoiceTime             string `excel:"主动语音总时长"`
	VoiceAvgTime          string `excel:"主动语音平均时长"`
	PassiveVoiceNum       string `excel:"被动语音次数"`
	PassiveVoiceRefuseNum string `excel:"被动语音拒接率"`
	PassiveVoiceRate      string `excel:"被动语音接听率"`
	PassiveVoiceTime      string `excel:"被动语音总时长"`
	PassiveVoiceAvgTime   string `excel:"被动语音平均时长"`
	VideoNum              string `excel:"主动视频次数"`
	VideoRefuseNum        string `excel:"主动视频被拒接次数"`
	VideoRate             string `excel:"主动视频接听率"`
	VideoTime             string `excel:"主动视频总时长"`
	VideoAvgTime          string `excel:"主动视频平均时长"`
	PassiveVideoNum       string `excel:"被动视频音次数"`
	PassiveVideoRefuseNum string `excel:"被动视频拒接率"`
	PassiveVideoRate      string `excel:"被动视频接听率"`
	PassiveVideoTime      string `excel:"被动视频总时长"`
	PassiveVideoAvgTime   string `excel:"被动视频平均时长"`
	SpeedVoiceNum         string `excel:"速配语音次数"`
	SpeedVoiceRefuseNum   string `excel:"速配语音被拒接次数"`
	SpeedVoiceRate        string `excel:"速配语音接听率"`
	SpeedVoiceTime        string `excel:"速配语音总时长"`
	SpeedVoiceAvgTime     string `excel:"速配语音平均时长"`
	SpeedVideoNum         string `excel:"速配视频次数"`
	SpeedVideoRefuseNum   string `excel:"速配视频被拒接次数"`
	SpeedVideoRate        string `excel:"速配视频接听率"`
	SpeedVideoTime        string `excel:"速配视频总时长"`
	SpeedVideoAvgTime     string `excel:"速配视频平均时长"`
	InvitationNum         string `excel:"邀请人数"` //invitation_binding这个表
	InvitationRechargeRmb string `excel:"邀请充值金额"`
	SayIncome             string `excel:"搭讪收益"`
	MsgIncome             string `excel:"消息收益"`
	VoiceIncome           string `excel:"语音收益"`
	SpeedVoiceIncome      string `excel:"速配语音收益"`
	VideoIncome           string `excel:"视频收益"`
	SpeedVideoIncome      string `excel:"速配视频收益"`
	InvitationIncome      string `excel:"邀请收益"`
	GiftIncome            string `excel:"礼物收益"`
	TotalIncome           string `excel:"总收益"`
}

func main() {
	invitationDataExcelList := make([]*InvitationDataExcel, 0)
	invitationDataExcelList = append(invitationDataExcelList, &InvitationDataExcel{Uid: "10", NickName: "zzy1"})
	invitationDataExcelList = append(invitationDataExcelList, &InvitationDataExcel{Uid: "20", NickName: "zzy2"})
	invitationDataExcelList = append(invitationDataExcelList, &InvitationDataExcel{Uid: "30", NickName: "zzy3"})
	file := excel.CreateExcel[InvitationDataExcel](invitationDataExcelList)
	file.SaveAs("./test.xlsx")
	//forwardInfoMaps := make(map[int]*ssh.ForwardInfo)
	//forwardInfoMaps[3306] = &ssh.ForwardInfo{
	//	Port: 3306,
	//	IP:   "db.panlian.com",
	//}
	//
	//go ssh.Conn("root", "120.78.197.32:9022", forwardInfoMaps, ssh.PublicKeyAuthFunc("./id_rsa"))
	//time.Sleep(time.Second * 3)
	//
	//chatLog, err := data.NewGormDB("127.0.0.1:3306", "readonly", "S1l,nv.9Gs#A", "chat_log", 1, 1, logger.Default.LogMode(logger.Info))
	//if err != nil {
	//	panic(err)
	//}
	//chatLog.Where("id = ?", "20")
	////_, err := data.NewGormDB("127.0.0.1:3306", "readonly", "S1l,nv.9Gs#A", "chat_log", 1, 1, logger.Default.LogMode(logger.Info))
	////if err != nil {
	////	panic(err)
	////}
	//
	//commCache := cache.NewDemoCommCache(nil)
	//local_cache.NewLocalCache[cache.DemoComm](time.Second*30, time.Second*30, commCache)

}
