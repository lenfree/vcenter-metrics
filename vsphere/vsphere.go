package vsphere

import (
        "net/url"
        "github.com/vmware/govmomi"
        "github.com/vmware/govmomi/find"
        "github.com/vmware/govmomi/object"
        "github.com/vmware/govmomi/property"
        "github.com/vmware/govmomi/vim25/mo"
        "github.com/vmware/govmomi/vim25/types"
        "golang.org/x/net/context"
        "strconv"
)

type Config struct {
        Host     string
        User     string
        Pass     string
        Insecure string
}

func (v Config) Connect(ctx context.Context) (*find.Finder, *govmomi.Client, *object.Datacenter, error){

        vcenterURL := "https://" + v.User + ":" + v.Pass + "@" + v.Host + "/sdk"

        u, _ := url.Parse(vcenterURL)

        c, err := govmomi.NewClient(ctx, u, v.insecure())
        if err != nil {
                return nil, nil, nil, err
        }

        f := find.NewFinder(c.Client, true)

        // Find one and only datacenter
        dc, err := f.DefaultDatacenter(ctx)
        if err != nil {
                return nil, nil, nil, err
        }

        // Make future calls local to this datacenter
        f.SetDatacenter(dc)
        return f, c, dc, nil
}

func QueryClusterVMs(ctx context.Context, c *govmomi.Client, f *find.Finder, dc *object.Datacenter, p string) ([]mo.VirtualMachine, error) {

        // Assume user pass empty string and default to current directory
        if len(p) == 0 {
                p = "*"
        }
        // Assume user specify which directory to search for VMs
        if p != "*" {
                p = dc.Common.InventoryPath + "/" + "vm" + "/" + p + "/" + "*"
        }

        cs, err := f.VirtualMachineList(ctx, p)
        if err != nil {
                return nil, err
        }
        pc := property.DefaultCollector(c.Client)

        var refs []types.ManagedObjectReference
        for _, c := range cs{
                refs = append(refs, c.Reference())
        }

        var vms []mo.VirtualMachine
        err = pc.Retrieve(ctx, refs, []string{"summary"}, &vms)
        if err != nil {
                return nil, err
        }

        return vms, nil
}

func QueryDatastore(ctx context.Context, c *govmomi.Client, f *find.Finder) ([]mo.Datastore, error) {
        // Find datastores in datacenter
        dss, err := f.DatastoreList(ctx, "*")
        if err != nil {
                return nil, err
        }

        pc := property.DefaultCollector(c.Client)

        var refs []types.ManagedObjectReference
        for _, ds := range dss {
                refs = append(refs, ds.Reference())
        }

        var datastores []mo.Datastore
        err = pc.Retrieve(ctx, refs, []string{"summary"}, &datastores)
        if err != nil {
                return nil, err
        }

        return datastores, nil
}

func QueryHosts(ctx context.Context, c *govmomi.Client, f *find.Finder) ([]mo.HostSystem, error) {
        cs, err := f.HostSystemList(ctx, "*")
        if err != nil {
                return nil, err
        }
        pc := property.DefaultCollector(c.Client)

        var refs []types.ManagedObjectReference
        for _, c := range cs{
                refs = append(refs, c.Reference())
        }

        var hosts []mo.HostSystem
        err = pc.Retrieve(ctx, refs, []string{"summary"}, &hosts)
        if err != nil {
                return nil, err
        }
        return hosts, nil
}

func (v Config) insecure() bool {
        if len(v.Insecure) > 0 {
                b, _ := strconv.ParseBool(v.Insecure)
                return b
        }
        return true
}

func (v Config) Logout(ctx context.Context) {
        v.Logout(ctx)
}
