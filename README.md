# ContainerScale

Note: This plugin is currently in development. Configuration, build patterns, etc may all change before this is stable.

A [CNI] plugin that connects containers to a Tailscale network. This means that as the container is created, they're added to your tsnet.

It is intended that you would use a Tailscale `AuthKey` to automatically authenticate these containers on your tsnet.

## Background

There are three popular networking abstractions for containers. [CNI], [Netavark], and Docker networking. The eventual goal of this repository is to support all three of these networking abstractions. For now, this repository only supports [CNI].

Different container runtimes support different networking abstractions:

- Kubernetes: [CNI]
- Podman: By default, [Netavark]. Can be configured to use [CNI] (for now).
- Docker: Docker networking.

## Setup & Configuration

More detailed instructions are TBA, but the gist is:

```bash
# Build the ContainerScale binary and drop it in a directory that CNI can find it.
go build -o /opt/cni/bin/containerscale
```

### Available Flags

The flags you can use to configure the plugin is:
```
AuthKey         Required  Authentication key from tailscale. 
TailscaledFlags Optional  Extra flags to run with `tailscaled`
TailscaleFlags  Optional  Extra flags to run with `tailscale up`
```

### Docker

Docker does not use [CNI]. Support for a docker network plugin is tracked in [#1].

### Podman
#### Netavark

[Netavark] is a new networking model that the Podman team is adopting as the default networking system. Support for a [Netavark] plugin is tracked in [#2].

#### CNI
First you need to make sure that Podman is running with CNI networking:

1. Open `/etc/containers/container.conf`. 
2. Find the line that starts with `#network_backend`. 
3. Uncomment it and change it to `network_backend = "cni"`.

Second, you need to create a new network configuration. You can do this at the user level.

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

## Kubernetes

This should just work with Kubernetes as well. I have not tried it yet. Work to improve documentation for kubernetes is tracked in [#3].

[#1]: https://github.com/aaomidi/ContainerScale/issues/1
[#2]: https://github.com/aaomidi/ContainerScale/issues/2
[#3]: https://github.com/aaomidi/ContainerScale/issues/3

[CNI]: https://github.com/containernetworking/cni
[Netavark]: https://github.com/containers/netavark