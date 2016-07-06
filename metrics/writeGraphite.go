package metrics

import (
        "fmt"
        "github.com/marpaia/graphite-golang"
)

type Config struct {
        Host string
        Port int
        Prefix string
}

func (m Config) SendMetric(metrics map[string]string) (error) {
        graphite, err := graphite.NewGraphite(m.Host, m.Port)
        if err != nil {
                fmt.Printf("Connection error with %s", err)
                return err
        }


        for k, v := range metrics {
                name := "vsphere" + "." + k
        // Send a stat
                err = graphite.SimpleSend(name, v)
                if err != nil {
                        fmt.Printf("Send Metrics error %s", err)
                        return err
                }
        }
        graphite.Disconnect()
        return nil
}

func (m Config) SendMetricName(n string, metrics map[string]string) (error) {
        graphite, err := graphite.NewGraphite(m.Host, m.Port)
        if err != nil {
                fmt.Printf("Connection error with %s", err)
                return err
        }


        for k, v := range metrics {
                name := m.Prefix + "." + n + "." + k
        // Send a stat
                err = graphite.SimpleSend(name, v)
                if err != nil {
                        fmt.Printf("Send Metrics error %s", err)
                        return err
                }
        }
        return nil
}
