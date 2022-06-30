package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Logstash struct {
	Hostname   string
	Port       int
	Connection *net.TCPConn
	Timeout    int
}

type LogEntry struct {
	LogType string `json:"logType"`
	Message string `json:"message"`
}

var Logger *Logstash

const INFO = "INFO"
const ERROR = "ERROR"
const WARNING = "WARNING"

func New(hostname string, port int, timeout int) *Logstash {
	l := Logstash{}
	l.Hostname = hostname
	l.Port = port
	l.Connection = nil
	l.Timeout = timeout
	return &l
}

func (l *Logstash) Dump() {
	fmt.Println("Hostname:   ", l.Hostname)
	fmt.Println("Port:       ", l.Port)
	fmt.Println("Connection: ", l.Connection)
	fmt.Println("Timeout:    ", l.Timeout)
}

func (l *Logstash) SetTimeouts() {
	deadline := time.Now().Add(time.Duration(l.Timeout) * time.Millisecond)
	l.Connection.SetDeadline(deadline)
	l.Connection.SetWriteDeadline(deadline)
	l.Connection.SetReadDeadline(deadline)
}

func (l *Logstash) Connect() (*net.TCPConn, error) {
	var connection *net.TCPConn
	service := fmt.Sprintf("%s:%d", l.Hostname, l.Port)
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return connection, err
	}
	connection, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return connection, err
	}
	if connection != nil {
		l.Connection = connection
		l.Connection.SetLinger(0) // default -1
		l.Connection.SetNoDelay(true)
		l.Connection.SetKeepAlive(true)
		l.Connection.SetKeepAlivePeriod(time.Duration(5) * time.Second)
		l.SetTimeouts()
	}
	return connection, err
}

func (l *Logstash) Writeln(message string) error {
	var err = errors.New("TCP Connection is nil.")
	message = fmt.Sprintf("%s\n", message)
	if l.Connection != nil {
		_, err = l.Connection.Write([]byte(message))
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				l.Connection.Close()
				l.Connection = nil
				if err != nil {
					return err
				}
			} else {
				l.Connection.Close()
				l.Connection = nil
				return err
			}
		} else {
			// Successful write! Let's extend the timeoul.
			l.SetTimeouts()
			return nil
		}
	}
	return err
}

func Logf(logType string, message string, v ...any) {
	message = fmt.Sprintf(message, v...)

	var logData, _ = json.Marshal(LogEntry{
		LogType: logType,
		Message: message,
	})

	var err = Logger.Writeln(string(logData))
	if err != nil {
		log.Printf("[LOGGER] Could not write data to LogStash. Error: %s\n", err.Error())

		if Logger.Connection == nil {
			log.Printf("[LOGGER] Retrying to connect to LogStash\n")
			_, err = Logger.Connect()
			if err != nil {
				log.Printf("[LOGGER] Could not reconnect to LogStash. Error: %s\n", err.Error())
			} else {
				log.Printf("[LOGGER] Successfully reconnected to LogStash\n")
			}
		}
	}

	log.Printf("[%s] %s\n",
		logType, message)
}
