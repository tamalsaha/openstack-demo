package main

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
     "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/pagination"
	"fmt"
	"os"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/mitchellh/mapstructure"

)

const auth_url  = "https://auth.cloud.ovh.net/"
func getAuth()gophercloud.AuthOptions  {
	return gophercloud.AuthOptions{
		IdentityEndpoint: auth_url+"v2.0",
		Username: os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		TenantID: os.Getenv("OS_TENANT_ID"),
		TenantName: os.Getenv("OS_TENANT_NAME"),

	}
}

func main()  {
	client, err := getOpenstacComputeClient()
	if err != nil {
		fmt.Println(err)
	}
	//createPublicNetwork(client)
	getServerAddresses(client)
	//getflavorlist(client)
	//createNetwork(client)
	//getNetworkList(client)
	//getServerList(client)
	/*getflavorlist(client)
	getNetworkList(client)
	getServerList(client)*/
	/*getSSHKey(client)
	getimagelist(client)

	createServer(client)*/

}

func getimagelist(client *gophercloud.ServiceClient)  error {
	opts := images.ListOpts{ChangesSince: "2014-01-01T01:02:03Z", Name: "Ubuntu 16.04"}
	pager := images.ListDetail(client, opts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}
		for _, i := range imageList {
			fmt.Println(i.ID)
		}
		return true, nil
	})

	return err
}


func getflavorlist(client *gophercloud.ServiceClient)  error {
	opts := flavors.ListOpts{ChangesSince: "2014-01-01T01:02:03Z"}
	pager := flavors.ListDetail(client, opts)

	return  pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		for _, f := range flavorList {
			fmt.Println(f)
			// "f" will be a flavors.Flavor
			fmt.Println(f.Name, "<-->", f.ID, f.RAM)
		}
		return true, nil
	})
}

func getServerAddresses(client *gophercloud.ServiceClient)  {
	server, err := servers.Get(client, "cbc97bf5-b8bd-4850-9f95-172af4077cde").Extract()
	fmt.Println(err)
	type Address struct {
		IpType string `mapstructure:"OS-EXT-IPS:type"`
		Addr   string
	}
	var addresses map[string][]Address
	err = mapstructure.Decode(server.Addresses, &addresses)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(server.AccessIPv4)
	for network, addrList := range addresses {
		for _, props := range addrList {
			fmt.Println(props.IpType, props.Addr)
		}
		fmt.Println(network)
	}
}

func getServerList(client *gophercloud.ServiceClient)  {
	opts := servers.ListOpts{}

	pager := servers.List(client, opts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		fmt.Println(err)
		for _, s := range serverList {
			fmt.Println(s.KeyName, s.Addresses, s.Status, s.AccessIPv4)
			fmt.Println(s.SecurityGroups)
			servers.ListAddresses(client, s.ID).EachPage(func(page pagination.Page) (bool, error) {
				as, err := servers.ExtractAddresses(page)
				fmt.Println(err, as)
				for b, a := range as {
					fmt.Println(b, a[0].Address, "<><><>")
				}
				return true, nil

			})
		}
		return true, nil
	})
	fmt.Println(err)
}

func createNetwork(client *gophercloud.ServiceClient) {
	tr := true
	net, err :=networks.Create(client, networks.CreateOpts{
		Name: "testnetwork",
		AdminStateUp: &tr,
	//	Shared: &tr,
		//TenantID: os.Getenv("OS_TENANT_ID"),
	}).Extract()
	fmt.Println(err, net.Status)
}

func getNetworkList(client *gophercloud.ServiceClient)  {
	tr := true
	opts := networks.ListOpts{Shared: &tr}

	// Retrieve a pager (i.e. a paginated collection)
	pager := networks.List(client, opts)

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		fmt.Println(err)
		for _, n := range networkList {
			// "n" will be a networks.Network
			fmt.Println(n.Name, n.ID, n.Status)
		}
		return true, nil
	})
	fmt.Println(err)
}

func createPublicNetwork(client *gophercloud.ServiceClient) {
	iTrue := true
	options := networks.CreateOpts{Name: "public", AdminStateUp: &iTrue, Shared: &iTrue, TenantID: "41b18bcc3f0e460ba97b19e5b6912247"}
	_, err := networks.Create(client, options).Extract()
	fmt.Println(err)
}

func createServer(client *gophercloud.ServiceClient)  {
	opts := servers.CreateOpts{
		Name:      "testovhserver",
		ImageRef:  "2e962277-13ad-44f1-9b0d-56e6b0ef1c00",
		FlavorRef: "550757b3-36c2-4027-b6fe-d70f45304b9c",
		Networks: []servers.Network{
			{
				UUID: "764d0ecb-f8a5-47d9-b034-53b5b61666a7",
			},
		},
		/*Personality: servers.Personality{
			{
				Path:     "/home/ubuntu/.ssh/authorized_keys",
				Contents: []byte(sshkey),
			},
		},*/
	}

	createOpts := keypairs.CreateOptsExt{
		CreateOptsBuilder: opts,
		KeyName: "sanjid",
	}

	server, err := servers.Create(client, createOpts).Extract()
	//servers.Update(client, server.ID, servers.UpdateOpts{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(server.Status, server.ID, server.KeyName)
}

func getSSHKey(client *gophercloud.ServiceClient)  {
	pager := keypairs.List(client)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		keylist, err := keypairs.ExtractKeyPairs(page)
		fmt.Println(err)
		for _, k := range keylist {
			fmt.Println(k.Name)
		}
		return true, nil
	})
}

func createSSHKey(client *gophercloud.ServiceClient)  {

}
func getOpenstacComputeClient() (*gophercloud.ServiceClient, error)  {
	provider, err := openstack.AuthenticatedClient(getAuth())
	if err != nil {
	return nil, err
	}


	return openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Name: "nova",
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

func getOpenstackNetworkClient() (*gophercloud.ServiceClient, error)  {
	provider, err := openstack.AuthenticatedClient(getAuth())
	fmt.Println(err)
	return openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name: "neutron",
		Region: os.Getenv("OS_REGION_NAME"),
	})
}
