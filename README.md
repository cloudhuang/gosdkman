SDKMAN

===========================

- [x] Setup the local JDK folder user - **.sdkman**
- [x] List all the local JDKs under the JDK folder
- [x] List all the JDKs remotely
- [x] Show the screen just like the sdkman
- [] jdk commands
    - [x] jdk list
    - [] jdk install VERSION
    - [] jdk use VERSION - change the JAVA_HOME to the VERSION

[] Update PATH environment - to add JAVA_HOME\bin item

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


