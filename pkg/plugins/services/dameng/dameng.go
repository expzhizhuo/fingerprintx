/*
  - Package dameng
    @Author: zhizhuo
    @IDE：GoLand
    @File: dameng.go
    @Date: 2025/2/12 下午4:44*
*/
package dameng

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/expzhizhuo/fingerprintx/pkg/plugins"
)

type DamengPlugin struct{}

const (
	DamengDB = "Dameng Database"
)

var (
	VerifyData = []byte{0x00, 0x00, 0x00, 0x00, 0x00}
)

func init() {
	plugins.RegisterPlugin(&DamengPlugin{})
}

func (p *DamengPlugin) Run(conn net.Conn, timeout time.Duration, target plugins.Target) (*plugins.Service, error) {
	// 构造链接数据包
	sendData := createBuildInfoCommand()
	// 发送数据包
	InfoResponse, err := sendDMCommand(conn, sendData, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send buildInfo command: %w", err)
	}
	// 解析buildInfo响应
	dmVersion, err := getDmVersion(InfoResponse)
	if err == nil {
		payload := plugins.ServiceDameng{}
		return plugins.CreateServiceFrom(target, payload, false,dmVersion , plugins.TCP), nil
	}

	return nil, fmt.Errorf("response did not match expected Dameng buildInfo response")
}

func (p *DamengPlugin) PortPriority(port uint16) bool {
	return port == 5236
}

func (p *DamengPlugin) Name() string {
	return DamengDB
}

func (p *DamengPlugin) Type() plugins.Protocol {
	return plugins.TCP
}

func (p *DamengPlugin) Priority() int {
	return 110
}

func createBuildInfoCommand() []byte {
	hexStr := "00000000c800520000000000000000000000009a000000000000000001020000000001090000000000000000000000000000000000000000000000000000000009000000382e312e322e3139320040000000182fa6db4e39692e3ad5559df6e0a026a77ce475f3097c784125089cd0f5c4de77d239d1946578dcf97840e514363ffdc71db4f7e1e22064e646006a7f7e19c1"
	cmd := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		decodedByte, _ := hex.DecodeString(hexStr[i : i+2])
		cmd[i/2] = decodedByte[0]
	}
	return cmd
}

func sendDMCommand(conn net.Conn, data []byte, timeout time.Duration) ([]byte, error) {
	// 设置连接的读写超时时间
	err := conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}
	// 发送数据到目标连接
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}
	// 接收完整的响应
	response, err := receiveFullResponse(conn, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}
	return response, nil
}

func receiveFullResponse(conn net.Conn, timeout time.Duration) ([]byte, error) {
	err := conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	var response []byte
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read from connection: %w", err)
		}
		response = append(response, buffer[:n]...)
		// 如果读取的数据小于缓冲区大小，可能是消息结束
		if n < len(buffer) {
			break
		}
	}
	return response, nil
}

func getDmVersion(data []byte) (string, error) {
	// 确认 VerifyData 是否存在于 data 中
	if !bytes.Contains(data, VerifyData) {
		return "", fmt.Errorf("data does not contain the required VerifyData pattern")
	}

	// 查找 '@' 字符的位置
	versionIndex := bytes.Index(data, []byte("@"))
	if versionIndex == -1 {
		return "", fmt.Errorf("version marker '@' not found")
	}

	// 确保 versionIndex 足够大，避免切片越界
	if versionIndex <= 12 {
		return "", fmt.Errorf("invalid version index: %d", versionIndex)
	}

	// 提取版本号
	version := data[versionIndex-12 : versionIndex-4]
	if len(version) == 0 {
		return "", fmt.Errorf("version data is empty")
	}

	// 将字节切片转换为字符串返回
	return string(version), nil
}