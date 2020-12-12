What's this?
===

Do you ever wanted to use your rabbit as a "network device"? The idea is to tunnel all ip packages via rabbit and using rabbit as a "vpn". For getting the ip frames I'm using Tuntap device and for shoveling golang.

It's a weekend hack and most probably not very stable. But I was able to have a stable ssh connection and a vnc connection

I do not have a sane use case, but it's working. It's the kind of projects "do stupid things, win stupid prices"

How to use/run/whatever
===
Either compile it (never checked if it's working on another machine) or download a pre-built binary.

For using it you need 2 hosts which can connect to a rabbitmq instance

On host1
```
sudo ./go-rabbit-tunnel "amqp://guest:guest@RABBITHOST:5672/" "10.1.0.10"
``` 

On host2
```
sudo ./go-rabbit-tunnel "amqp://guest:guest@RABBITHOST:5672/" "10.1.0.11"
```

On host1 Try ping
```
ping 10.1.0.11
```

Now you can try to ssh into the second host. Maybe it works on your machine :)
