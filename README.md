# ddregister

This code provides Dynamic DNS Registration for any hosts within your
domain, which is registered through Google Domains. The code should be run as a daemon,
as it checks the IP address of your NAT domain using [ipify](https://www.ipify.org/) and
re-registers with Google if your IP address has changed. You can set the interval
for the registration. To use,

```
bash% go get -u github.com/djthorpe/ddregister
bash% ddregister -help
Usage of ddregister:
  -debug
    	Set debugging mode
  -dns.passwd string
    	DNS registry password
  -dns.user string
    	DNS registry username
  -host string
    	Hostname to register
  -interval duration
    	IP address discovery interval (default 1h0m0s)
  -log.append
    	When writing log to file, append output to end of file
  -log.file string
    	File for logging (default: log to stderr)
  -verbose
    	Verbose logging
```

The `-dns.user` and `-dns.passwd` flags are required and can be retrieved as per
[these intructions](https://support.google.com/domains/answer/6147083?hl=en-GB).
For example,

```
bash% ddregister.go \
  -dns.user   XXXXXXXXXXXXXXX \
  -dns.passwd XXXXXXXXXXXXXXX \
  -host       XXXXXXXXXXXXXXX.mydomain.com \
  -interval 15s
```

Will run the code in the foreground and re-register your dynamic DNS entry
every 15 minutes.

