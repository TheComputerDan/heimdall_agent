package comms

import (
	"context"
	ds "github.com/TheComputerDan/heimdall/proto/dockerService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/TheComputerDan/heimdall/docker_agent/docker/connect"
	"github.com/TheComputerDan/heimdall/docker_agent/host"

	"net"
)

type server struct{}

func Start() {
	listener, err := net.Listen("tcp", ":8096")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	ds.RegisterDockerAgentServer(srv, &server{})
	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *server) Containers(request *ds.ContainersRequest, stream ds.DockerAgent_ContainersServer) error {

	hostContainers := connect.Containers()

	for _, container := range hostContainers {
		var Ports []*ds.Port
		var NetworkSettings *ds.SummaryNetworkSettings
		var Mounts []*ds.MountPoint

		for _, mount := range container.Mounts {
			mp := &ds.MountPoint{
				Destination: mount.Destination,
				Driver:      mount.Driver,
				Mode:        mount.Mode,
				Name:        mount.Name,
				Propagation: string(mount.Propagation),
				RW:          mount.RW,
				Source:      mount.Source,
				Type:        string(mount.Type),
				//UNSURE if these typecasts will work properly
			}
			Mounts = append(Mounts, mp)
		}

		for k, netSetting := range container.NetworkSettings.Networks {

			ipamConfig := &ds.EndpointIPAMConfig{
				IPv4Address:  netSetting.IPAMConfig.IPv4Address,
				IPv6Address:  netSetting.IPAMConfig.IPv6Address,
				LinkLocalIPs: netSetting.IPAMConfig.LinkLocalIPs,
			}

			network := make(map[string]*ds.EndpointSettings)
			network[k] = &ds.EndpointSettings{
				Gateway:             netSetting.Gateway,
				IPAMConfig:          ipamConfig,
				Links:               netSetting.Links,
				Aliases:             netSetting.Aliases,
				NetworkID:           netSetting.NetworkID,
				EndpointID:          netSetting.EndpointID,
				IPAddress:           netSetting.IPAddress,
				IPPrefixLen:         int64(netSetting.IPPrefixLen),
				IPv6Gateway:         netSetting.IPv6Gateway,
				GlobalIPv6Address:   netSetting.GlobalIPv6Address,
				GlobalIPv6PrefixLen: int64(netSetting.GlobalIPv6PrefixLen),
				MacAddress:          netSetting.MacAddress,
				DriverOpts:          netSetting.DriverOpts,
			}

			NetworkSettings = &ds.SummaryNetworkSettings{
				Networks: network,
			}
		}

		for _, port := range container.Ports {
			p := &ds.Port{
				IP:          port.IP,
				PrivatePort: uint32(port.PrivatePort),
				PublicPort:  uint32(port.PublicPort),
			}
			Ports = append(Ports, p)
		}

		resp := &ds.ContainersResponse{
			ID:              container.ID,
			Names:           container.Names,
			State:           container.State,
			Labels:          container.Labels,
			Status:          container.Status,
			Ports:           Ports,
			Command:         container.Command,
			Created:         container.Created,
			ImageID:         container.ImageID,
			Image:           container.Image,
			SizeRw:          container.SizeRw,
			SizeRootFs:      container.SizeRootFs,
			NetworkSettings: NetworkSettings,
			Mounts:          Mounts,
			NetworkMode:     container.HostConfig.NetworkMode,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Images(request *ds.ImagesRequest, stream ds.DockerAgent_ImagesServer) error {
	hostImages := connect.Images()

	for _, image := range hostImages {

		resp := &ds.ImagesResponse{
			Containers:  image.Containers,
			Created:     image.Created,
			ID:          image.ID,
			Labels:      image.Labels,
			ParentID:    image.ParentID,
			RepoDigests: image.RepoDigests,
			RepoTags:    image.RepoTags,
			SharedSize:  image.SharedSize,
			Size:        image.Size,
			VirtualSize: image.VirtualSize,
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Host(ctx context.Context, request *ds.HostRequest) (*ds.HostResponse, error) {
	hostInfo := host.BuildInfo()

	return &ds.HostResponse{
		Hostname: hostInfo.Hostname,
		OsType:   hostInfo.OSType,
		Ipv4:     hostInfo.IP["ipv4"],
		Ipv6:     hostInfo.IP["ipv6"],
	}, nil
}
