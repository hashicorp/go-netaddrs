## Command Line Tool Usage

Inspired by [go-discover](https://github.com/hashicorp/go-discover), `go-netaddrs` is a Go (golang) library and command line tool to discover ip addresses of nodes in a customizable fashion suitable for any environment.

Install the command line tool with:

```
go get -u github.com/hashicorp/go-netaddrs/cmd/netaddrs
```

Example usage

```
$ netaddrs -q ip "exec=/usr/local/bin/query_ec2.sh"

# Output
172.25.19.221 172.25.24.182 172.25.21.52
```

```
$ netaddrs -q ip "exec=discover -q addrs provider=aws region=us-west-2 tag_key=consul-server tag_value=true"

# Output
172.25.19.221 172.25.24.182 172.25.21.52
```

```
$ netaddrs -q ip "consul-cluster.private.consul.11eb5b95-2882-0215-b2c7-0242ac11000d.aws.hcp.dev"

# Output
172.25.19.221 172.25.24.182 172.25.21.52
```

## Library Usage

Install the library with:

```
go get -u github.com/hashicorp/go-netaddrs
```

Usage sample:
```
import netaddrs "github.com/hashicorp/go-netaddrs"

func ServerAddresses(server_addresses_cfg string, logger hclog.Logger) ([]string, error) {


   // Example `server_addresses_cfg` values:
   // consul-cluster.private.consul.11eb5b95-2882-0215-b2c7-0242ac11000d.aws.hcp.dev
   // exec=query_ec2.sh
   // exec=discover -q addrs provider=aws region=us-west-2 tag_key=consul-server tag_value=true

   addresses, err := netaddrs.IPAddrs(server_addresses_cfg, logger)
   if err != nil {
       logger.Error("Error retrieving server addresses", err)
       return nil, err
   }


   logger.Info("Server addresses", addresses)
   return addresses, err
}
```

## Testing

```bash
$ go test
```