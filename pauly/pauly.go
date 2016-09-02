package pauly

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Fields map[string]interface{}

type PaulyLogger struct {
	Connection  net.Conn
	Environment string
	Application string
	Host        string
	Port        int
}

func New(environment, application, host string, port int) PaulyLogger {
	conn, err := net.Dial("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Print("Could not connect to logger", err)
	}

	return PaulyLogger{
		Connection:  conn,
		Environment: environment,
		Application: application,
		Host:        host,
		Port:        port,
	}
}

func serialize(fields Fields, application, severity string) (string, error) {
	fields["severity"] = severity
	fields["application"] = application
	data, err := json.Marshal(fields)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return string(data), nil
}

func (pauly *PaulyLogger) Debug(fields Fields) {
	data, _ := serialize(fields, pauly.Application, "DEBUG")
	if pauly.Environment == "development" {
		log.Print(data)
	}
	if pauly.Environment == "production" {
		fmt.Fprintf(pauly.Connection, data)
	}
}

func (pauly *PaulyLogger) Info(fields Fields) {
	data, _ := serialize(fields, pauly.Application, "INFO")
	if pauly.Environment == "development" {
		log.Print(data)
	}
	if pauly.Environment == "production" {
		fmt.Fprintf(pauly.Connection, data)
	}
}

func (pauly *PaulyLogger) Warn(fields Fields) {
	data, _ := serialize(fields, pauly.Application, "WARN")
	if pauly.Environment == "development" {
		log.Print(data)
	}
	if pauly.Environment == "production" {
		fmt.Fprintf(pauly.Connection, data)
	}
}

func (pauly *PaulyLogger) Error(fields Fields) {
	data, _ := serialize(fields, pauly.Application, "ERROR")
	if pauly.Environment == "development" {
		log.Print(data)
	}
	if pauly.Environment == "production" {
		fmt.Fprintf(pauly.Connection, data)
	}
}

func (pauly *PaulyLogger) Fatal(fields Fields) {
	data, _ := serialize(fields, pauly.Application, "FATAL")
	if pauly.Environment == "development" {
		log.Print(data)
	}
	if pauly.Environment == "production" {
		fmt.Fprintf(pauly.Connection, data)
	}
	log.Fatal(data)
}

func (pauly *PaulyLogger) Close() {
	pauly.Connection.Close()
}
