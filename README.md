# ContainerScale

A [CNI] plugin that connects containers to a Tailscale network. This means that as the container is created, they're added to your tsnet.

It is intended that you would use a Tailscale `AuthKey` to automatically authenticate these containers on your tsnet.

## Installation

More detailed instructions are TBA, but the gist is:

```bash
go build -o /opt/cni/bin/containerscale
```

## Configuration

### Flags

```
AuthKey         Required  Authentication key from tailscale. 
TailscaledFlags Optional  Extra flags to run with `tailscaled`
TailscaleFlags  Optional  Extra flags to run with `tailscale up`
```

### Podman
You will need to create a new `conflist` network configuration. You can do this at the user level.

```bash
touch ~/.config/cni/net.d/99-containerscale.conflist
```

Example Configuration:

```json
{
  "cniVersion": "0.4.0",
  "name": "myts",
  "plugins": [
    {
      "type": "bridge",
      "bridge": "cni-podman1",
      "isGateway": true,
      "ipMasq": true,
      "hairpinMode": true,
      "ipam": {
        "type": "host-local",
        "routes": [
          {
            "dst": "0.0.0.0/0"
          }
        ],
        "ranges": [
          [
            {
              "subnet": "10.89.0.0/24",
              "gateway": "10.89.0.1"
            }
          ]
        ]
      },
      "capabilities": {
        "ips": true
      }
    },
    {
      "type": "containerscale",
      "runtimeConfig": {
        "authKey": "tskey-auth-#####",
        "tailscaleFlags": [
          "--ssh" 
        ]
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    },
    {
      "type": "firewall",
      "backend": ""
    },
    {
      "type": "tuning"
    }
  ]
}
```

[CNI]: https://github.com/containernetworking/cni