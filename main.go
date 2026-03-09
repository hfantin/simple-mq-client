package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("use simple-mq-cli publish|consume [queueName] [message]")
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("arquivo .env não encontrado")
		os.Exit(1)
	}

	queueName := os.Args[2]
	hostPort := os.Getenv("MQ_HOST")
	manager := os.Getenv("MQ_MANAGER")
	channel := os.Getenv("MQ_CHANNEL")

	fmt.Println("conectando em ", hostPort, "gerenciador", manager, "canal", channel)

	// Connect options
	cd := ibmmq.NewMQCD()
	cd.ChannelName = channel
	cd.ConnectionName = hostPort

	csp := ibmmq.NewMQCSP()
	// csp.AuthenticationType = ibmmq.MQCSP_AUTH_NONE
	co := ibmmq.NewMQCNO()
	co.ClientConn = cd
	co.SecurityParms = csp

	qMgr, err := ibmmq.Connx(manager, co)
	if err != nil {
		fmt.Println("erro de conexao:", err)
		return
	}
	defer qMgr.Disc()

	od := ibmmq.NewMQOD()
	od.ObjectName = queueName
	openOptions := ibmmq.MQOO_OUTPUT | ibmmq.MQOO_INPUT_AS_Q_DEF

	queue, err := qMgr.Open(od, openOptions)
	if err != nil {
		fmt.Println("erro ao abrir fila:", err)
		return
	}
	defer queue.Close(0)

	switch os.Args[1] {
	case "publish":
		publicarMensagem(&queue, &queueName)
	case "consume":
		consumirMensagem(&queue, &queueName)
	default:
		fmt.Println("comando desconhecido. utilize publish ou consume.")
	}
}

func publicarMensagem(queue *ibmmq.MQObject, queueName *string) {
	if len(os.Args) < 4 {
		fmt.Println("use simple-mq-cli publish [queueName] [message]")
		return
	}
	msg := os.Args[3]
	fmt.Println("publicando mensagem ", msg, "na fila", *queueName)
	md := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()
	err := queue.Put(md, pmo, []byte(msg))
	if err != nil {
		fmt.Println("erro na publicacao:", err)
	} else {
		fmt.Println("mensagem publicada!")
	}
}

// func consumirMensagem(queue *ibmmq.MQObject, queueName *string) {
// 	fmt.Println("consumindo mensagem da fila", *queueName)
// 	md := ibmmq.NewMQMD()
// 	gmo := ibmmq.NewMQGMO()
// 	buf := make([]byte, 1024)
// 	n, err := queue.Get(md, gmo, buf)
// 	if err != nil {
// 		fmt.Println("erro no consumo:", err)
// 	} else {
// 		fmt.Printf("recebido: %s\n", string(buf[:n]))
// 	}
// }

func consumirMensagem(queue *ibmmq.MQObject, queueName *string) {
	fmt.Println("consumindo mensagem da fila", *queueName)
	md := ibmmq.NewMQMD()
	gmo := ibmmq.NewMQGMO()
	buf := make([]byte, 64*1024) // 64KB buffer
	n, err := queue.Get(md, gmo, buf)
	if err != nil {
		if mqret, ok := err.(*ibmmq.MQReturn); ok {
			switch mqret.MQRC {
			case ibmmq.MQRC_TRUNCATED_MSG_FAILED:
				fmt.Println("mensagem muito grande para o buffer")
			case ibmmq.MQRC_NO_MSG_AVAILABLE:
				fmt.Println("nenhuma mensagem encontrada")
			default:
				fmt.Println("erro no consumo codigo", mqret.MQRC)
			}
		} else {
			fmt.Println("erro no consumo:", err)
		}
	} else {
		body := extractBody(buf[:n])
		// fmt.Printf("body: %s\n", body)
		// fmt.Printf("recebido: %s\n", string(buf[:n]))
		fmt.Printf("recebido: %s\n", body)
	}
}

func extractBody(msg []byte) []byte {
	// Step 1: Skip MQRFH2 header if present
	if len(msg) > 8 && string(msg[:4]) == "RFH " {
		headerLen := int(binary.BigEndian.Uint32(msg[4:8]))
		if len(msg) >= headerLen {
			msg = msg[headerLen:]
		}
	}
	// Step 2: Skip JMS folders (look for "</jms>")
	endTag := []byte("</jms>")
	idx := bytes.Index(msg, endTag)
	if idx != -1 && len(msg) > idx+len(endTag) {
		// Skip past the end of </jms>
		msg = msg[idx+len(endTag):]
		// Optionally trim leading whitespace
		msg = bytes.TrimLeft(msg, " \t\r\n")
	}
	return msg
}
