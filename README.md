JDKMAN

===========================

- [x] Setup the local JDK folder user - **.jdkman**
- [x] List all the local JDKs under the JDK folder
- [x] List all the JDKs remotely
- [x] Show the screen just like the sdkman
- [] jdk commands
    - [x] jdk list
    - [] jdk install VERSION
    - [] jdk use VERSION - change the JAVA_HOME to the VERSION

[] Update PATH environment - to add JAVA_HOME\bin item

### System environments cache
The environment variables are cached when a process starts.  Unless the process itself changes them they will not be visible until the process restarts.  In your case the batch file (running in a separate process) will update the env vars but the main process won't be able to see the changes.  There isn't a workaround short of making the same changes in the main process.  However if all you want to do is confirm that the env vars were changed then you can run another process that confirms the env vars were changed properly.  All new processes would get the new env vars.