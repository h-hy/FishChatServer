package provider

import (
    "errors"
    "bytes"
)

type LBSInfo struct {
    Cid    string
    Lac    string
    Rssi    string
    Arfcn    string
    Bsic    string
    Mnc    string
    Mcc    string
}
type GPSInfo struct {
    N_S_value    string
    N_S    string
    E_W_value    string
    E_W    string
    Speed    string
}

type Location struct {
    Type    byte
    LBSInfo []LBSInfo
    GPSInfo GPSInfo
}
// L|3571,9763,00,460,50|100|100
// G|100008.000,A,2232.4679,N,11356.7805,E,0.204,89.22,210911,,|7.49|152.6|100|100

func (self *Location) IsGPS() bool {
    return self.Type=='G'
}
func (self *Location) IsLBS() bool {
    return self.Type=='L'
}
func (self *Location) Parse(locatonRaw string) error {
    var locationDataSlice []string
    //开始按|分割
    locationData:=[]byte(locatonRaw)
    var data_spilt []byte  = []byte{'|'}
    nowindex:=0
    last_index:=len(locationData)
    for {
        index := bytes.Index(locationData[nowindex:],data_spilt)
        if (index==-1){
            if (nowindex != last_index){
                arg :=string(locationData[nowindex:last_index])
                locationDataSlice = append(locationDataSlice, arg)
            }
            break;
        }
        arg :=string(locationData[nowindex:nowindex+index])
        locationDataSlice = append(locationDataSlice, arg)
        nowindex+=index+1
    }
    //按|分割完成：locationDataSlice
    if (locationDataSlice[0]=="G"){
        self.Type='G'
        self.GPSInfo.Parse(locationDataSlice)
    }else if (locationDataSlice[0]=="L"){
        self.Type='L'
        self.LBSInfoParse(locationDataSlice)
    }
    return nil
}

// L   15653,31040,33,64,49,0,460:11032,31040,29,72,39,0,460:15651,31040,24,79,10,0,460:11033,31040,18,80,28,0,460:40402,26952,17,67,15,0,460:11031,31040,15,70,58,0,460:40732,31040,11,68,63,0,460,50   100   90
func (self *Location) LBSInfoParse(lbsRaw []string) error {
    var locationDataSlice []string
    locationData:=[]byte(lbsRaw[1])
    var data_spilt []byte  = []byte{':'}
    nowindex:=0
    last_index:=len(locationData)
    for {
        index := bytes.Index(locationData[nowindex:],data_spilt)
        if (index==-1){
            if (nowindex != last_index){
                arg :=string(locationData[nowindex:last_index])
                locationDataSlice = append(locationDataSlice, arg)
            }
            break;
        }
        arg :=string(locationData[nowindex:nowindex+index])
        locationDataSlice = append(locationDataSlice, arg)
        nowindex+=index+1
    }
    for i := 0; i < len(locationDataSlice); i++ {
        LBSCell,err :=self.LoadLBSCell(locationDataSlice[i])
        if (err==nil){
            self.LBSInfo=append(self.LBSInfo,LBSCell)
        }
    }

    return nil
}
// 15653,31040,33,64,49,0,460
func (self *Location) LoadLBSCell(cell string) (LBSInfo, error) {
    var LBSCellSlice []string
    locationData:=[]byte(cell)
    var data_spilt []byte  = []byte{','}
    nowindex:=0
    last_index:=len(locationData)
    for {
        index := bytes.Index(locationData[nowindex:],data_spilt)
        if (index==-1){
            if (nowindex != last_index){
                arg :=string(locationData[nowindex:last_index])
                LBSCellSlice = append(LBSCellSlice, arg)
            }
            break;
        }
        arg :=string(locationData[nowindex:nowindex+index])
        LBSCellSlice = append(LBSCellSlice, arg)
        nowindex+=index+1
    }
    var LbSCell LBSInfo
    if (len(LBSCellSlice)<7){
        return LbSCell,errors.New("Cell Error")
    }
    LbSCell.Cid=LBSCellSlice[0]
    LbSCell.Lac=LBSCellSlice[1]
    LbSCell.Rssi=LBSCellSlice[2]
    LbSCell.Arfcn=LBSCellSlice[3]
    LbSCell.Bsic=LBSCellSlice[4]
    LbSCell.Mnc=LBSCellSlice[5]
    LbSCell.Mcc=LBSCellSlice[6]
    return LbSCell,nil
}

// G|100008.000,A,2232.4679,N,11356.7805,E,0.204,89.22,210911,,|7.49|152.6|100|100
func (self *GPSInfo) Parse(gpsRaw []string) error {
    var GPSDataSlice []string
    locationData:=[]byte(gpsRaw[1])
    var data_spilt []byte  = []byte{','}
    nowindex:=0
    last_index:=len(locationData)
    for {
        index := bytes.Index(locationData[nowindex:],data_spilt)
        if (index==-1){
            if (nowindex != last_index){
                arg :=string(locationData[nowindex:last_index])
                GPSDataSlice = append(GPSDataSlice, arg)
            }
            break;
        }
        arg :=string(locationData[nowindex:nowindex+index])
        GPSDataSlice = append(GPSDataSlice, arg)
        nowindex+=index+1
    }
    if (len(GPSDataSlice)<7){
        return errors.New("GPSInfo Error")
    }
    self.N_S_value=GPSDataSlice[2]
    self.N_S=GPSDataSlice[3]
    self.E_W_value=GPSDataSlice[4]
    self.E_W=GPSDataSlice[5]
    self.Speed=GPSDataSlice[6]

    return nil
}