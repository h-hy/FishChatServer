package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type LBSInfo struct {
	Cid   string
	Lac   string
	Rssi  string
	Arfcn string
	Bsic  string
	Mnc   string
	Mcc   string
}
type GPSInfo struct {
	N_S_value string
	N_S       string
	E_W_value string
	E_W       string
	Speed     string
}

type Location struct {
	Type         byte
	LBSInfo      []*LBSInfo
	GPSInfo      GPSInfo
	Energy       int
	LocationData LocationDataStruct
}

// L|3571,9763,00,460,50|100|100
// G|100008.000,A,2232.4679,N,11356.7805,E,0.204,89.22,210911,,|7.49|152.6|100|100

func (self *Location) IsGPS() bool {
	return self.Type == 'G'
}
func (self *Location) IsLBS() bool {
	return self.Type == 'L'
}
func (self *Location) Parse(locatonRaw string) error {
	var locationDataSlice []string
	//开始按|分割
	locationData := []byte(locatonRaw)
	var data_spilt []byte = []byte{'|'}
	nowindex := 0
	last_index := len(locationData)
	for {
		index := bytes.Index(locationData[nowindex:], data_spilt)
		if index == -1 {
			if nowindex != last_index {
				arg := string(locationData[nowindex:last_index])
				locationDataSlice = append(locationDataSlice, arg)
			}
			break
		}
		arg := string(locationData[nowindex : nowindex+index])
		locationDataSlice = append(locationDataSlice, arg)
		nowindex += index + 1
	}
	//按|分割完成：locationDataSlice
	if locationDataSlice[0] == "G" {
		self.Type = 'G'
		self.GPSInfo.Parse(locationDataSlice)
	} else if locationDataSlice[0] == "L" {
		self.Type = 'L'
		self.LBSInfoParse(locationDataSlice)
		locationReturn, err := self.LoadLocationInfo()
		if err != nil {
			return nil
		}
		self.LocationData = (*locationReturn).Result
	}
	energy, err := strconv.Atoi(locationDataSlice[len(locationDataSlice)-1])
	if err == nil {
		self.Energy = energy
	}
	return nil
}

// L   15653,31040,33,64,49,0,460:11032,31040,29,72,39,0,460:15651,31040,24,79,10,0,460:11033,31040,18,80,28,0,460:40402,26952,17,67,15,0,460:11031,31040,15,70,58,0,460:40732,31040,11,68,63,0,460,50   100   90
func (self *Location) LBSInfoParse(lbsRaw []string) error {
	var locationDataSlice []string
	locationData := []byte(lbsRaw[1])
	var data_spilt []byte = []byte{':'}
	nowindex := 0
	last_index := len(locationData)
	for {
		index := bytes.Index(locationData[nowindex:], data_spilt)
		if index == -1 {
			if nowindex != last_index {
				arg := string(locationData[nowindex:last_index])
				locationDataSlice = append(locationDataSlice, arg)
			}
			break
		}
		arg := string(locationData[nowindex : nowindex+index])
		locationDataSlice = append(locationDataSlice, arg)
		nowindex += index + 1
	}
	for i := 0; i < len(locationDataSlice); i++ {
		LBSCell, err := self.LoadLBSCell(locationDataSlice[i])
		if err == nil {
			self.LBSInfo = append(self.LBSInfo, LBSCell)
		}
	}

	return nil
}

// 15653,31040,33,64,49,0,460
func (self *Location) LoadLBSCell(cell string) (*LBSInfo, error) {
	var LBSCellSlice []string
	locationData := []byte(cell)
	var data_spilt []byte = []byte{','}
	nowindex := 0
	last_index := len(locationData)
	for {
		index := bytes.Index(locationData[nowindex:], data_spilt)
		if index == -1 {
			if nowindex != last_index {
				arg := string(locationData[nowindex:last_index])
				LBSCellSlice = append(LBSCellSlice, arg)
			}
			break
		}
		arg := string(locationData[nowindex : nowindex+index])
		LBSCellSlice = append(LBSCellSlice, arg)
		nowindex += index + 1
	}
	LbSCell := &LBSInfo{}
	if len(LBSCellSlice) < 7 {
		return LbSCell, errors.New("Cell Error")
	}
	LbSCell.Cid = LBSCellSlice[0]
	LbSCell.Lac = LBSCellSlice[1]
	LbSCell.Rssi = LBSCellSlice[2]
	LbSCell.Arfcn = LBSCellSlice[3]
	LbSCell.Bsic = LBSCellSlice[4]
	LbSCell.Mnc = LBSCellSlice[5]
	LbSCell.Mcc = LBSCellSlice[6]
	return LbSCell, nil
}

// Cid    string
// Lac    string
// Rssi    string
// Arfcn    string
// Bsic    string
// Mnc    string
// Mcc    string
func (self *Location) LoadLocationInfo() (*locationReturn, error) {
	locationReturn := &locationReturn{}
	if self.Type != 'L' {
		return locationReturn, errors.New("Not LBS Location")
	}
	if len(self.LBSInfo) < 1 {
		return locationReturn, errors.New("LBSInfo Empty")
	}
	//开始构建url
	bts := "bts=" + self.LBSInfo[0].Mcc + "," + self.LBSInfo[0].Mnc + "," + self.LBSInfo[0].Lac + "," + self.LBSInfo[0].Cid + "," + self.LBSInfo[0].Rssi
	var nearbts []string
	for i := 1; i < len(self.LBSInfo); i++ {
		nearbts = append(nearbts, self.LBSInfo[0].Mcc+","+self.LBSInfo[0].Mnc+","+self.LBSInfo[0].Lac+","+self.LBSInfo[0].Cid+","+self.LBSInfo[0].Rssi)
	}
	client := &http.Client{}

	//开始发出请求
	resp, err := client.Get("http://apilocate.amap.com/position?key=698028f2d74a36a25c5af5bea759b482&accesstype=0&imei=352315052834187&cdma=0&bts=" + bts + "&nearbts=" + strings.Join(nearbts, "|"))

	defer resp.Body.Close()

	//开始读取回复
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return locationReturn, err
	}
	//开始解析json
	err = json.Unmarshal(body, locationReturn) // JSON to Struct
	return locationReturn, nil
}

type LocationDataStruct struct {
	Type     string `json:"type"`
	Location string `json:"location"`
	Eadius   string `json:"radius"`
	Desc     string `json:"desc"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Citycode string `json:"citycode"`
	Adcode   string `json:"adcode"`
	Road     string `json:"road"`
	Street   string `json:"street"`
	Poi      string `json:"poi"`
}
type locationReturn struct {
	Status   uint64 `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Result   LocationDataStruct
}

// G|100008.000,A,2232.4679,N,11356.7805,E,0.204,89.22,210911,,|7.49|152.6|100|100
func (self *GPSInfo) Parse(gpsRaw []string) error {
	var GPSDataSlice []string
	locationData := []byte(gpsRaw[1])
	var data_spilt []byte = []byte{','}
	nowindex := 0
	last_index := len(locationData)
	for {
		index := bytes.Index(locationData[nowindex:], data_spilt)
		if index == -1 {
			if nowindex != last_index {
				arg := string(locationData[nowindex:last_index])
				GPSDataSlice = append(GPSDataSlice, arg)
			}
			break
		}
		arg := string(locationData[nowindex : nowindex+index])
		GPSDataSlice = append(GPSDataSlice, arg)
		nowindex += index + 1
	}
	if len(GPSDataSlice) < 7 {
		return errors.New("GPSInfo Error")
	}
	self.N_S_value = GPSDataSlice[2]
	self.N_S = GPSDataSlice[3]
	self.E_W_value = GPSDataSlice[4]
	self.E_W = GPSDataSlice[5]
	self.Speed = GPSDataSlice[6]

	return nil
}
