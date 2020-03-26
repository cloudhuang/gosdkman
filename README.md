SDKMAN
===========================

GoSDKMan is a command-line tool which allows you to easily install, manage, and work with multiple Java environments for windows system.

> If you already installed the JDK by the executable file, please remove the corresponding env like `C:\ProgramData\Oracle\Java\javapath` from your **PATH**.

## Usage

### Help
```
> sdk.exe -h
Usage: GoSDKMan [options...]

GoSDKMan is a command-line tool which allows you to easily install, manage, and work with multiple Java environments for windows.

Options:
  -h, --help            Display usage
  -v, --version         Display version
  -i, --install         Install the new JDK version
  -d, --uninstall       Uninstall the JDK version
  -u, --use             Use the installed JDK version
  -l, --list            List all the available versions
```

### List all available JDK
```
> sdk.exe -l
=====================================================================================
Available Java Versions
=====================================================================================
 Vendor          | Use | Version      | Dist       | Status     | Identifier
-------------------------------------------------------------------------------------
 Amazon          |     | 1.8          | amaz       |            | 1.8-amaz
                 |     | 11           | amaz       |            | 11-amaz
 AdoptOpenJDK    |     | 8u242+b08    | adpt       |            | 8u242+b08-adpt
                 |     | 9.0.4+11     | adpt       |            | 9.0.4+11-adpt
                 |     | 10.0.2+13.1  | adpt       |            | 10.0.2+13.1-adpt
                 |     | 11.0.6+10    | adpt       |            | 11.0.6+10-adpt
                 |     | 12.0.2+10    | adpt       |            | 12.0.2+10-adpt
                 |     | 13.0.2+8     | adpt       |            | 13.0.2+8-adpt
                 |     | 14+36        | adpt       |            | 14+36-adpt
 Java.net        |     | 1.9_181      | open       |            | 1.9_181-open
                 |     | 11_28        | open       |            | 11_28-open
                 |     | 12_32        | open       |            | 12_32-open
                 |     | 13_33        | open       |            | 13_33-open
                 |     | 14_36        | open       |            | 14_36-open
                 |     | 1.7.0_75-b13 | open       |            | 1.7.0_75-b13-open
                 |     | 1.8.0_41-b04 | open       |            | 1.8.0_41-b04-open
 Azul Zulu       | >>> | 13.29.9      | zulu       | Installed  | 13.29.9-zulu
                 |     | 14.27.1      | zulu       | Installed  | 14.27.1-zulu
=====================================================================================
Use the Identifier for installation:
         sdk -i 11.0.6.10.1-amaz
```

### Install new JDK version
```
> sdk -i  11-amaz

Downloading amazon-corretto-11-x64-windows-jdk.zip [Total Size: 185 MB]... 41.70%

......

Downloading amazon-corretto-11-x64-windows-jdk.zip [Total Size: 185 MB]... 100%
2020/03/26 15:37:28
Download completed in 19m21.5544776s
The new JDK version installed success, please restart the console.

```

then start a new console
```
> java -version
openjdk version "11.0.6" 2020-01-14 LTS
OpenJDK Runtime Environment Corretto-11.0.6.10.1 (build 11.0.6+10-LTS)
OpenJDK 64-Bit Server VM Corretto-11.0.6.10.1 (build 11.0.6+10-LTS, mixed mode)
```


## 解决Github访问问题
由于众所周知的问题，在国内访问github不是太顺畅，比如`raw.githubusercontent.com`这个地址就无法访问，需要通过修改hosts的方式来解决这个问题:

- 打开Dns检测|[Dns查询 – 站长工具](http://tool.chinaz.com/dns?type=1&host=raw.githubusercontent.com&ip=)
- 在检测输入栏中输入GitHub官网，比如`raw.githubusercontent.com`
- 把检测列表里某个IP配置到hosts里，并对应写上github官网域名
```
151.101.228.133  raw.githubusercontent.com
```


## System environments cache
The environment variables are cached when a process starts.  Unless the process itself changes them they will not be visible until the process restarts.  In your case the batch file (running in a separate process) will update the env vars but the main process won't be able to see the changes.  There isn't a workaround short of making the same changes in the main process.  However if all you want to do is confirm that the env vars were changed then you can run another process that confirms the env vars were changed properly.  All new processes would get the new env vars.


