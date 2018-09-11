# SOAP-GO-OSC
Open Shop Channel's SOAP Server, designed specifically to handle Wii Shop Channel SOAP.

## I wanna fork just because I can!
Please refrain from forking if you don't plan on contributing, tinkering or working with it.
It only just makes unnecessary messiness.

## What's the difference between this repo and that other SOAP repo?
This is the SOAP Server Software. The other repository only has the communication templates between a Wii and WSC's server.

## How do i configure the server?
This the the content of the configuration file, the **config.xml**:
```xml
<config>
	<!-- 2018 is the default port, 80 the standard port of the WSC. -->
	<setting key="port" value="2018" /> 
</config>
```
You can edit the port at the `value="<port>"`.
