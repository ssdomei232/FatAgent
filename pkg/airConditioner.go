package pkg

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
)

// Modbus 寄存器地址常量
const (
	coilStartAddress     = 1
	coilCount            = 32
	registerStartAddress = 3
	registerCount        = 3
	temperatureRegister  = 6 // 40007
	humidityRegister     = 7 // 40008
)

type ACStatus struct {
	// 开关量 (BOOL)
	FanProtect        bool `json:"fanProtect"`        // 送风机保护
	FanPressureSwitch bool `json:"fanPressureSwitch"` // 送风压差开关
	PhaseProtect      bool `json:"phaseProtect"`      // 相许保护
	FireLinkage       bool `json:"fireLinkage"`       // 防火连锁开关
	HeaterOverheat    bool `json:"heaterOverheat"`    // 电加热过热
	RemoteSwitch      bool `json:"remoteSwitch"`      // 远程开关
	HumidifierFault   bool `json:"humidifierFault"`   // 加湿器故障
	Comp1HighPressure bool `json:"comp1HighPressure"` // 1#压机高压
	Comp1LowPressure  bool `json:"comp1LowPressure"`  // 1#压机低压
	Comp1Overload     bool `json:"comp1Overload"`     // 1#压机过载
	Cond1Overload     bool `json:"cond1Overload"`     // 1#冷凝风机过载
	Comp2HighPressure bool `json:"comp2HighPressure"` // 2#压机高压
	Comp2LowPressure  bool `json:"comp2LowPressure"`  // 2#压机低压
	Comp2Overload     bool `json:"comp2Overload"`     // 2#压机过载
	Cond2Overload     bool `json:"cond2Overload"`     // 2#冷凝风机过载
	RunSignal         bool `json:"runSignal"`         // 运行信号
	FaultSignal       bool `json:"faultSignal"`       // 故障信号
	Valve             bool `json:"valve"`             // 风阀
	SupplyFan         bool `json:"supplyFan"`         // 送风机
	Comp1             bool `json:"comp1"`             // 1#压机
	Cond1             bool `json:"cond1"`             // 1#冷凝风机
	Comp2             bool `json:"comp2"`             // 2#压机
	Cond2             bool `json:"cond2"`             // 2#冷凝风机
	Heater1           bool `json:"heater1"`           // 1#电加热
	Heater2           bool `json:"heater2"`           // 2#电加热

	// 模拟量 (WORD)
	SupplyTemp     uint16 // 送风温度
	IndoorTemp     uint16 // 室内温度
	IndoorHumidity uint16 // 室内湿度
}

// 控制器
type ACController struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func NewACController(port string) (*ACController, error) {
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 2 * time.Second

	if err := handler.Connect(); err != nil {
		handler.Close() // 避免资源泄漏
		return nil, fmt.Errorf("连接Modbus设备失败: %v", err)
	}

	client := modbus.NewClient(handler)

	return &ACController{
		handler: handler,
		client:  client,
	}, nil
}

// 启动
func (ac *ACController) Start() error {
	_, err := ac.client.WriteSingleCoil(0, 0xFF00) // 启动 -> 40001
	if err != nil {
		return fmt.Errorf("启动设备失败: %v", err)
	}
	return nil
}

// 停止
func (ac *ACController) Stop() error {
	_, err := ac.client.WriteSingleCoil(1, 0xFF00) // 停止 -> 40002
	if err != nil {
		return fmt.Errorf("停止设备失败: %v", err)
	}
	return nil
}

// 设置温度
func (ac *ACController) SetTemperature(value uint16) error {
	_, err := ac.client.WriteSingleRegister(temperatureRegister, value) // 40007
	if err != nil {
		return fmt.Errorf("设置温度失败: %v", err)
	}
	return nil
}

// 设置湿度
func (ac *ACController) SetHumidity(value uint16) error {
	_, err := ac.client.WriteSingleRegister(humidityRegister, value) // 40008
	if err != nil {
		return fmt.Errorf("设置湿度失败: %v", err)
	}
	return nil
}

func (ac *ACController) GetAllStatus() (*ACStatus, error) {
	status := &ACStatus{}

	// 1. 读取布尔量
	bools, err := ac.client.ReadCoils(coilStartAddress, coilCount) // offset=1 -> 40002
	if err != nil {
		return nil, fmt.Errorf("读取布尔量失败: %v", err)
	}
	if len(bools)*8 < 25 {
		return nil, fmt.Errorf("布尔量数据不足，缺少必要位")
	}
	status.FanProtect = getBool(bools, 0)
	status.FanPressureSwitch = getBool(bools, 1)
	status.PhaseProtect = getBool(bools, 2)
	status.FireLinkage = getBool(bools, 3)
	status.HeaterOverheat = getBool(bools, 4)
	status.RemoteSwitch = getBool(bools, 5)
	status.HumidifierFault = getBool(bools, 6)
	status.Comp1HighPressure = getBool(bools, 7)
	status.Comp1LowPressure = getBool(bools, 8)
	status.Comp1Overload = getBool(bools, 9)
	status.Cond1Overload = getBool(bools, 10)
	status.Comp2HighPressure = getBool(bools, 11)
	status.Comp2LowPressure = getBool(bools, 12)
	status.Comp2Overload = getBool(bools, 13)
	status.Cond2Overload = getBool(bools, 14)
	status.RunSignal = getBool(bools, 15)
	status.FaultSignal = getBool(bools, 16)
	status.Valve = getBool(bools, 17)
	status.SupplyFan = getBool(bools, 18)
	status.Comp1 = getBool(bools, 19)
	status.Cond1 = getBool(bools, 20)
	status.Comp2 = getBool(bools, 21)
	status.Cond2 = getBool(bools, 22)
	status.Heater1 = getBool(bools, 23)
	status.Heater2 = getBool(bools, 24)

	// 2. 读取测量量：送风温度、室内温度、室内湿度
	// 寄存器：40004, 40005, 40006 → offset=3, count=3
	words, err := ac.client.ReadHoldingRegisters(registerStartAddress, registerCount)
	if err != nil {
		return nil, fmt.Errorf("读取模拟量失败: %v", err)
	}
	if len(words) < 6 {
		return nil, fmt.Errorf("模拟量数据不足，期望6字节，实际%d字节", len(words))
	}

	status.SupplyTemp = binary.BigEndian.Uint16(words[0:2])
	status.IndoorTemp = binary.BigEndian.Uint16(words[2:4])
	status.IndoorHumidity = binary.BigEndian.Uint16(words[4:6])

	return status, nil
}

// 解析单个 bit
func getBool(data []byte, index int) bool {
	byteIndex := index / 8
	bitIndex := index % 8
	if byteIndex >= len(data) {
		return false
	}
	return (data[byteIndex] & (1 << bitIndex)) != 0
}
