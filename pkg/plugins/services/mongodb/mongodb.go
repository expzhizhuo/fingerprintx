package mongodb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/expzhizhuo/fingerprintx/pkg/plugins"
)

type MongoDBPlugin struct{}

const (
	MONGODB = "MongoDB"
)

func init() {
	plugins.RegisterPlugin(&MongoDBPlugin{})
}

func (p *MongoDBPlugin) Run(conn net.Conn, timeout time.Duration, target plugins.Target) (*plugins.Service, error) {
	// 构造buildInfo命令
	buildInfoCmd := createBuildInfoCommand()
	// 发送buildInfo命令
	buildInfoResponse, err := sendMongoDBCommand(conn, buildInfoCmd, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send buildInfo command: %w", err)
	}
	// 解析buildInfo响应
	buildInfoDoc, err := parseBSONResponse(buildInfoResponse)
	if err == nil {
		payload := plugins.ServiceMongoDB{
			PacketType:   "buildInfo",
			ErrorMessage: "",
			ErrorCode:    0,
		}
		return plugins.CreateServiceFrom(target, payload, false, buildInfoDoc["version"].(string), plugins.TCP), nil
	}
	return nil, fmt.Errorf("response did not match expected MongoDB buildInfo response")
}

func (p *MongoDBPlugin) PortPriority(port uint16) bool {
	return port == 27017
}

func (p *MongoDBPlugin) Name() string {
	return MONGODB
}

func (p *MongoDBPlugin) Type() plugins.Protocol {
	return plugins.TCP
}

func (p *MongoDBPlugin) Priority() int {
	return 2
}

func createBuildInfoCommand() []byte {
	hexStr := "3b0000000100000000000000d40700000000000061646d696e2e24636d640000000000ffffffff14000000106275696c64496e666f000100000000"
	cmd := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		decodedByte, _ := hex.DecodeString(hexStr[i : i+2])
		cmd[i/2] = decodedByte[0]
	}
	return cmd
}

func sendMongoDBCommand(conn net.Conn, cmd []byte, timeout time.Duration) ([]byte, error) {
	err := conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

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

	// 读取消息头
	header := make([]byte, 16)
	_, err = conn.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read response header: %w", err)
	}

	var messageLength int32
	buf := bytes.NewReader(header[:4])
	err = binary.Read(buf, binary.LittleEndian, &messageLength)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	remainingBytes := int(messageLength) - 16
	body := make([]byte, remainingBytes)
	_, err = conn.Read(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	response := append(header, body...)
	return response, nil
}

func parseBSONResponse(response []byte) (map[string]interface{}, error) {
	if len(response) < 36 {
		return nil, fmt.Errorf("response is too short")
	}
	// 跳过前36个字节的消息头，直接解析BSON数据
	messageBody := response[36:]
	var doc map[string]interface{}
	err := bson.Unmarshal(messageBody, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal BSON response: %w", err)
	}
	return doc, nil
}
