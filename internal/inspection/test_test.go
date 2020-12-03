package inspection

import (
	"context"
	"testing"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/spf13/viper"
)

func TestInspection(t *testing.T) {
	v := viper.New()
	v.Set("logging.console.level", "debug")
	v.Set("host.roles", "ch")

	cfg := config.NewSettings(v)

	schema := schema.GetClickhouseCheckSchema(cfg)

	for _, r := range schema.Reports {
		r.Gather(context.Background())
	}

	println("======================")
	RunInspection(cfg)
	println("======================")
}

/*
   cpu := &reports.CPUReport{}
   cpu.Gather(context.Background())
   fmt.Printf("cpu: %+v\n", cpu)

   disk := &reports.DiskReport{}
   disk.Gather(context.Background())
   fmt.Printf("disk: %+v\n", disk)

   drv := &reports.DriverReport{
       DriverName: "intel_sgx",
   }
   drv.Gather(context.Background())
   fmt.Printf("drv intel_sgx: %+v\n", drv)

   drv2 := &reports.DriverReport{
       DriverName: "isgx",
   }
   drv2.Gather(context.Background())
   fmt.Printf("drv isgx: %+v\n", drv2)

   docker := &reports.DockerReport{}
   docker.Gather(context.Background())
   fmt.Printf("docker: %+v\n", docker)

   http := &reports.HTTPReport{
       URL: "https://registry.aggregion.com",
   }
   http.Gather(context.Background())
   fmt.Printf("http: %+v\n", http)

   hv := &reports.HypervisorReport{}
   hv.Gather(context.Background())
   fmt.Printf("hv: %+v\n", hv)

   kern := &reports.KernelReport{}
   kern.Gather(context.Background())
   fmt.Printf("kern: %+v\n", kern)

   net := &reports.NetProbeReport{
       Type:   "tcp",
       Target: "185.175.44.42:18766",
   }
   net.Gather(context.Background())
   fmt.Printf("net: %+v\n", net)

   os := &reports.OSReport{}
   os.Gather(context.Background())
   fmt.Printf("os: %+v\n", os)

   ram := &reports.RAMReport{}
   ram.Gather(context.Background())
   fmt.Printf("ram: %+v\n", ram)

   se := &reports.OSSeLinuxReport{}
   se.Gather(context.Background())
   fmt.Printf("se: %+v\n", se)

   swap := &reports.SwapFileReport{}
   swap.Gather(context.Background())
   fmt.Printf("swap: %+v\n", swap)
*/
