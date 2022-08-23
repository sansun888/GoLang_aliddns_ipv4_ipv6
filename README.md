# GoLang_aliddns_ipv4_ipv6
动态域名解析设置,自动把本机的公共地址设置到阿里解析服务器上

# 
有看到别人的python版本,但是windows上需要单独安装python有点麻烦<br>
这个go 可以编译成exe直接执行

# windows
go build -o aliddns.exe ./main

# linux
go build -o aliddns ./main

# 运行反馈效果
```
C:\Users\tick>D:\ddns\GoLang_aliddns_ipv4_ipv6_Win_X64\aliddns.exe
[LOG] {"AccessKeyId":"***","AccessSecret":"***","Ipv4Flag":1,"Ipv6Flag":1,"Domain":"baise.tk","NameIpv4":"t41",
"NameIpv6":"t61","LogFileFlag":1,"Ip4Url":"https://api-ipv4.ip.sb/ip","Ip6Url":"https://api6.ipify.org"}
取得ip地址:112.19.***
添加解析: t41.baise.tk
添加ip4成功
取得ip地址:2409:8a62:5b1d:6970:***
添加解析: t61.baise.tk
添加ip6成功

C:\Users\tick>
C:\Users\tick>
C:\Users\tick>D:\ddns\GoLang_aliddns_ipv4_ipv6_Win_X64\aliddns.exe
[LOG] {"AccessKeyId":"***","AccessSecret":"***","Ipv4Flag":1,"Ipv6Flag":1,"Domain":"baise.tk","NameIpv4":"t41",
"NameIpv6":"t61","LogFileFlag":1,"Ip4Url":"https://api-ipv4.ip.sb/ip","Ip6Url":"https://api6.ipify.org"}
取得ip地址:112.19.***
ip地址没变
取得ip地址:2409:8a62:5b1d:6970:***
ip地址没变

```
