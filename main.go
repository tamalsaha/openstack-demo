package main

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
     "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/pagination"
	"fmt"
	"os"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

const auth_url  = "https://auth.cloud.ovh.net/"
func getAuth()gophercloud.AuthOptions  {
	return gophercloud.AuthOptions{
		IdentityEndpoint: auth_url+"v2.0",
		Username: os.Getenv("OS_USER_NAME"),
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
	getflavorlist(client)
	getServerList(client)
	getimagelist(client)
	getNetworkList(client)
	createServer(client)

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
	opts := flavors.ListOpts{ChangesSince: "2014-01-01T01:02:03Z", MinRAM: 1}
	pager := flavors.ListDetail(client, opts)
	return  pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		for _, f := range flavorList {
			// "f" will be a flavors.Flavor
			fmt.Println(f.Name, "<-->", f.ID)
		}
		return true, nil
	})
}

func getServerList(client *gophercloud.ServiceClient)  {
	opts := servers.ListOpts{}

	pager := servers.List(client, opts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		fmt.Println(err)
		for _, s := range serverList {
			fmt.Println(s.Name)
		}
		return true, nil
	})
	fmt.Println(err)
}

func getNetworkList(client *gophercloud.ServiceClient)  {
	opts := networks.ListOpts{}

	// Retrieve a pager (i.e. a paginated collection)
	pager := networks.List(client, opts)

	// Define an anonymous function to be executed on each page's iteration
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		fmt.Println(err)
		for _, n := range networkList {
			// "n" will be a networks.Network
			fmt.Println(n.Name, n.ID)
		}
		return true, nil
	})
	fmt.Println(err)
}

func createServer(client *gophercloud.ServiceClient)  {

	server, err := servers.Create(client, servers.CreateOpts{
		Name: "testovhserver",
		ImageRef: "2e962277-13ad-44f1-9b0d-56e6b0ef1c00",
		FlavorRef: "550757b3-36c2-4027-b6fe-d70f45304b9c",
		Networks: []servers.Network{
			{
				UUID: "4ef3de54-af73-4e57-a205-5e1d97c48eb1",
			},
		},


	}).Extract()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(server.Status, server.ID)
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

