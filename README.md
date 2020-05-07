# J2735 Decode Agent

Go package for decoding J2735 CV2X standards. C source is compiled using ASN.1 compiler (https://github.com/vlm/asn1c). CGO is used for wrapping C calls from Go. Client interface uses MQTT to consume messages and decodes in either XML or JSON format for consumption. 

For build instructions, please see Dockerfile.

## How this works

The decode agent create bridges between two MQTT brokers and either operates in batch or streaming mode. It can convert the binary blob provided by the broker it subscribe to into XML or JSON. It can also passthrough the messages to the publishing broker. 

The decode agent provides a HTTP interface in which PUT calls can be made to modify the MQTT settings (more on this in the following section). Finally, the decode agent can operate in playback mode which can be used to replay simply logs encoded in hex strings, please see [cmd/decode-agent/test/resources/logs/](https://gitlab.verizon.com/hsuyu/j2735-decode-agent/-/tree/master/cmd/decode-agent/test/resources/logs) for more details.

## Configuration

Only one yaml file is needed to create multiple decode agents. Examples can be seem in the compose/config/config.yaml file. Each yaml section needs to separated by 3 hyphens as per yaml specs. The yaml schema are described below:

```
name: string    # name of the decode agent, the http endpoint will use this name
subscribe:      # settings for connecting to subscribing broker
    clientid: string    # clientid for mqtt broker to subscribe to
    server: string      # mqtt broker to subscribe to
    topic: string       # topic to subscribe to
    qos: number         # quality of service 
publish:        # settings for connecting to publishing broker
    clientid: string    # clientid for mqtt broker to publish to
    server: string      # mqtt broker to subscribe to
    topic: string       # topic to subscribe to
    qos: number         # quality of service 
op:             # operation settings
    mode: string        # mode to operate in => stream/batch/playback
    format: string      # format to publish in => json/xml/pass 
    http-auth: string   # path to the location of the http password
    batchconfig:        # config for batch mode operations
        pubfreq: number     # how often the batched message are published
    playbackconfig:
        file: string        # file to read from 
        loop: boolean       # whether to restart playback file once it ends
        pubfreq: number     # how often to publish the message read
---
.... second endpoint

```

A couple of notes:
*  the name represents the http endpoint (e.g. name: streamAgent => <fqdn>.com/streamAgent/settings)
*  all fields in the yaml configuration can be overriden with an environment variables, for example to change operating mode: 

      op.mode => set `<NAME>`_OP_MODE where NAME is the name of the decode agent
*  the http-auth file must follow the folowing text format:

```
username
password
```

### Command Line Flags

*  Log levels (trace=-1, debug=0, info=1, warn=2, error=3) for the program can be changed using the **-loglevel=<-1,0,1,2,3,4,5>**.
*  Configuration file locations (default location is /etc/decode-agent/config.yaml) for the program can be specified using **-cfg=<file_path>**.  

## Components

The decode agent is composed of the following components:

*  pkg/decoder => The CGO components which wraps a c library compiled using ASN.1 compiler.
*  internal/cfgparser => The config parsing component which defines the structure of the configuration yaml file.
*  cmd/decode-agent => The cmd line application that handles bridging the MQTT brokers.

## Testing

Unit tests are included each of the above components (e.g.):

`go test -v .`

Docker compose is setup to run a smoke test with a local broker (e.g.):

`docker-compose up`
