# tunl.online

<b>tunl</b>.online is a developer tool for easy make your service available on the internet and easy way to inspect incoming traffic with view of request Headers and Body.

NOTE: current version in alpha mode.

### Download

From [release page](https://github.com/black40x/tunl-cli/releases)

### Build

```
go run build.go
```

### Examples

Share local port:

```
tunl http 8000
```

Share IP address:

```
tunl http 192.168.1.10:8000
```

Example output:

```
🚀 tunl started!

Version              0.1.0-Alpha
Session expired at   2022-08-19 17:09:53
Web monitor          http://127.0.0.1:6060
Forwarding           http://127.0.0.1:8000 -> http://10-10-10-10-abcde.tunl.online

HTTP Requests:  

[5.9ms] [GET] /
...
```

![](assets/web-mon.png)