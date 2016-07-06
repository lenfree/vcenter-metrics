package main

import (
        "fmt"
        "os"
        "strconv"
        "golang.org/x/net/context"
        "github.com/vmware/govmomi"
        "github.com/vmware/govmomi/find"
        "github.com/pivotal-golang/bytefmt"
        "git02.ae.sda.corp.telstra.com/vcenter-metrics/pauly"
        // Remove line 13
        "git02.ae.sda.corp.telstra.com/vcenter-metrics/metrics"
        "git02.ae.sda.corp.telstra.com/vcenter-metrics/vsphere"
        "github.com/vmware/govmomi/object"
        "strings"
        "sync"
        "time"
)

var logger pauly.PaulyLogger
// Remove line 22
var graphite metrics.Config
var v vsphere.Config
var interval Interval

type Interval struct {
        t int
}

func init() {
        //do not use proxy
        os.Unsetenv("http_proxy")
        os.Unsetenv("https_proxy")
        i, _ := strconv.Atoi(os.Getenv("QUERY_INTERVAL"))
        environment := os.Getenv("ENVIRONMENT")
        loggerHost := os.Getenv("LOGGER_HOST")
        loggerPort, _ := strconv.Atoi(os.Getenv("LOGGER_PORT"))
        vcenterHost := os.Getenv("VCENTER_HOST")
        vcenterUser := os.Getenv("VCENTER_USERNAME")
        vcenterPass := os.Getenv("VCENTER_PASSWORD")
        vcenterInsecure := os.Getenv("VCENTER_INSECURE")
        // Remove line 42 and 43
        metricHost := os.Getenv("METRIC_HOST")
        metricPort, _ := strconv.Atoi(os.Getenv("METRIC_PORT"))
        interval = Interval{
                t: i,
        }
        v = vsphere.Config{
                Host:     vcenterHost,
                User:     vcenterUser,
                Pass:     vcenterPass,
                Insecure: vcenterInsecure,
        }
        // Remove line 54 to 58
        graphite = metrics.Config{
                Host: metricHost,
                Port: metricPort,
                Prefix: "vcenter",
        }
        logger = pauly.New(
                        environment,
                        "vcenter-metrics",
                        loggerHost,
                        loggerPort)
}

func exit(err error) {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        logger.Error(pauly.Fields{
                "message": err.Error(),
        })
        os.Exit(1)
}

func main() {
        // Query Cluster Summary more frequently
        for {
                var wg sync.WaitGroup

                ctx, cancel := context.WithCancel(context.Background())
                f, c, dc, err := v.Connect(ctx)

                if err != nil {
                        exit(err)
                }

                logger.Info(pauly.Fields{
                        "message": "Login to vSphere...",
                })

                wg.Add(4)
                defer cancel()

                go queryDS(ctx, c, f, &wg)
                go vmCount(ctx, c, f, dc, &wg)
                go queryHosts(ctx, c, f, &wg)
                go vmsSummary(ctx, c, f, dc, &wg)

                wg.Wait()

                c.Logout(ctx)
                logger.Info(pauly.Fields{
                        "message": "Log out from vSphere...",
                })

                time.Sleep(intToTime(interval) * time.Minute)
                logger.Info(pauly.Fields{
                        "message": "Sleep for " + strconv.Itoa(interval.t) + " minutes...",
                })

        }
}

func queryDS(ctx context.Context, c *govmomi.Client, f *find.Finder, wg *sync.WaitGroup) {
        dsts, err := vsphere.QueryDatastore(ctx, c, f)
        if err != nil {
                logger.Error(pauly.Fields{
                        "message":  err.Error(),
                })
        }

        for _, dst := range dsts {
                logger.Info(pauly.Fields{
                        "message":        "datastore summary",
                        "datastore_name": whitespaceReplace(dst.Summary.Name),
                        "capacity_mb":    toMB(dst.Summary.Capacity),
                        "free_mb":        toMB(dst.Summary.FreeSpace),
                        "used_mb":        toMB(dst.Summary.Capacity - dst.Summary.FreeSpace),
                })

                // Remove line 133 to line 150
                graphite.Prefix = "vsphere" + "." + "datastore"
                metric := map[string]string{}
                metric["capacity"] = int64ToString(toMB(dst.Summary.Capacity))
                metric["free"] = int64ToString(toMB(dst.Summary.FreeSpace))
                metric["percent_free"] = strconv.Itoa(percentage(toFloat(dst.Summary.FreeSpace), toFloat(dst.Summary.Capacity)))
                metric["percent_used"] = strconv.Itoa(percentage(toFloat(used(dst.Summary.Capacity, dst.Summary.FreeSpace)),
                                                 toFloat(toMB(dst.Summary.Capacity))))

                err := graphite.SendMetricName(whitespaceReplace(dst.Summary.Name), metric)
                if err == nil {
                        logger.Info(pauly.Fields{
                                "message":  "Send metrics success",
                        })
                } else {
                        logger.Error(pauly.Fields{
                                "message":  err.Error(),
                        })
                }
        }
        wg.Done()
}

func vmCount(ctx context.Context, c *govmomi.Client, f *find.Finder, dc *object.Datacenter, wg *sync.WaitGroup) {
        vms, err := vsphere.QueryClusterVMs(ctx, c, f, dc, "*")
        if err != nil {
                logger.Error(pauly.Fields{
                        "message":  err.Error(),
                })
        }
        logger.Info(pauly.Fields{
                "message":        "total vm count",
                "count":          len(vms),
        })

        // Remove line 155 to 181
        graphite.Prefix = "vsphere"
        metric := map[string]string{}
        metric["total_count"] = int64ToString(int64(len(vms)))

        err = graphite.SendMetricName("vm", metric)
        if err == nil {
                logger.Info(pauly.Fields{
                        "message":  "Send metrics success",
                })
        } else {
                logger.Error(pauly.Fields{
                        "message":  err.Error(),
                })
        }
        wg.Done()
}

func vmsSummary(ctx context.Context, c *govmomi.Client, f *find.Finder, dc *object.Datacenter, wg *sync.WaitGroup) {
        // root path VMs
        vms, err := vsphere.QueryClusterVMs(ctx, c, f, dc, "*")
        if err != nil {
                logger.Error(pauly.Fields{
                        "message":  err.Error(),
                })
        }

        // Big-IP VMs in DO_NOT_DELETE FOLDER directory
        bigIPs, err := vsphere.QueryClusterVMs(ctx, c, f, dc, "DO_NOT_DELETE")
        if err != nil {
                logger.Error(pauly.Fields{
                        "message":  err.Error(),
                })
        }

        // Combine slices of []mo.VirtualMachine
        vms = append(vms, bigIPs...)

        for _, vm := range vms {
                logger.Info(pauly.Fields{
                        "message":                "vm summary",
                        "GuestToolsRunningStatus": vm.Summary.Guest.ToolsRunningStatus,
                        "Template":                vm.Summary.Config.Template,
                        "VmPathName":              vm.Summary.Config.VmPathName,
                        "Name":                    vm.Summary.Config.Name,
                        "NumVirtualDisks":         vm.Summary.Config.NumVirtualDisks,
                        "NumCpu":                  vm.Summary.Config.NumCpu,
                        "UptimeSeconds":           vm.Summary.QuickStats.UptimeSeconds,
                        "OverallCpuUsage":         vm.Summary.QuickStats.OverallCpuUsage,
                        "OverallCpuDemand":        vm.Summary.QuickStats.OverallCpuDemand,
                        "GuestMemoryUsage":        vm.Summary.QuickStats.GuestMemoryUsage,
                        "HostMemoryUsage":         vm.Summary.QuickStats.HostMemoryUsage,
                        "MemorySizeMB":            vm.Summary.Config.MemorySizeMB,
                        "OverallStatus":           vm.Summary.OverallStatus,
                        "Annotation":              vm.Summary.Config.Annotation,
                        "ConsumedOverheadMemory":  vm.Summary.QuickStats.ConsumedOverheadMemory,
                        "SwappedMemory":           vm.Summary.QuickStats.SwappedMemory,
                        "PowerState":              vm.Summary.Runtime.PowerState,
                        "ConnectionState":         vm.Summary.Runtime.ConnectionState,
                })
        }
        wg.Done()
}

func queryHosts(ctx context.Context, c *govmomi.Client, f *find.Finder, wg *sync.WaitGroup) {
        hosts, _ := vsphere.QueryHosts(ctx, c, f)
        for _, h := range hosts {
                mF := uint64(toMB(h.Summary.Hardware.MemorySize)) - uint64(h.Summary.QuickStats.OverallMemoryUsage)
                cF := totalCPU(h.Summary.Hardware.CpuMhz, h.Summary.Hardware.NumCpuCores) - int64(h.Summary.QuickStats.OverallCpuUsage)

                logger.Info(pauly.Fields{
                        "message":         "ESX host summary",
                        "name":            whitespaceReplace(h.Summary.Config.Name),
                        "memory_total_mb": toMB(h.Summary.Hardware.MemorySize),
                        "memory_used_mb":  h.Summary.QuickStats.OverallMemoryUsage,
                        "memory_free_mb":  mF,
                        "cpu_total_mhz":   totalCPU(h.Summary.Hardware.CpuMhz, h.Summary.Hardware.NumCpuCores),
                        "cpu_used_mhz":    h.Summary.QuickStats.OverallCpuUsage,
                        "cpu_free_mhz":    cF,
                })

                graphite.Prefix = "vsphere" + "." + "esx_host"
                metric := map[string]string{}
                metric["memory_total"] = int64ToString(toMB(h.Summary.Hardware.MemorySize))
                metric["memory_used"] = int32ToString(h.Summary.QuickStats.OverallMemoryUsage)
                metric["memory_free"] = uint64ToString(mF)
                metric["memory_percent_used"] = strconv.Itoa(percentage(int32ToFloat(h.Summary.QuickStats.OverallMemoryUsage),
                                                 toFloat(toMB(h.Summary.Hardware.MemorySize))))
                metric["memory_percent_free"] = strconv.Itoa(percentage(float64(mF), toFloat(toMB(h.Summary.Hardware.MemorySize))))
                metric["cpu_total"] = int64ToString(totalCPU(h.Summary.Hardware.CpuMhz, h.Summary.Hardware.NumCpuCores))
                metric["cpu_used"] = int32ToString(h.Summary.QuickStats.OverallCpuUsage)
                metric["cpu_free"] = int64ToString(cF)
                metric["cpu_percent_used"] = strconv.Itoa(percentage(int32ToFloat(h.Summary.QuickStats.OverallCpuUsage),
                                                 toFloat(totalCPU(h.Summary.Hardware.CpuMhz, h.Summary.Hardware.NumCpuCores))))
                metric["cpu_percent_free"] = strconv.Itoa(percentage(float64(cF), toFloat(totalCPU(h.Summary.Hardware.CpuMhz, h.Summary.Hardware.NumCpuCores))))

                err := graphite.SendMetricName(whitespaceReplace(h.Summary.Config.Name), metric)
                if err == nil {
                        logger.Info(pauly.Fields{
                                "message":  "Send metrics success",
                        })
                } else {
                        logger.Error(pauly.Fields{
                                "message":  err.Error(),
                        })
                }
      }
        wg.Done()
}

func totalCPU(h int32, c int16) int64 {
        return int64(h) * int64(c)
}

func toFloat(i int64) float64 {
        return float64(i)
}

func int32ToFloat(i int32) float64 {
        return float64(i)
}

func toMB(s int64) int64 {
        return s / 1024 / 1024
}

func int64ToString(i int64) string {
        return strconv.FormatInt(i, 10)
}

func uint64ToString(i uint64) string {
        return strconv.FormatInt(int64(i), 10)
}

func int32ToString(i int32) string {
        return strconv.FormatInt(int64(i), 10)
}

func intToString(i int32) string {
        return strconv.FormatInt(int64(i), 10)
}

func toPercent(i float64) int {
        return round(i * 100)
}

func round(val float64) int {
        if val < 0 { return int(val-0.5) }
        return int(val+0.5)
}

func used(t, u int64) int64 {
        return toMB(t) - toMB(u)
}

// Remove this block since this is not used
func toByteSize(b uint64) string {
        return bytefmt.ByteSize(uint64(b))
}

func whitespaceReplace(n string) string {
        return strings.Replace(n, " ", "_", -1)
}

func percentage(d, t float64) int {
        return toPercent(d / t)
}

func intToTime(i Interval) time.Duration {
        if i.t > 0 {
                return time.Duration(i.t)
        }
        return time.Duration(5)
}

// Remove this block not used
func ghzToString(i int64) string {
        return toMhzString(i)
}

// Remove this block not used
func toMhzString(i int64) string {
        return strconv.FormatInt(i, 10) + "Mhz"
}
