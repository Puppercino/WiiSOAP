# WiiSOAP
A SOAP Server slowly in development, designed specifically to handle Wii Shop Channel SOAP

## What's the difference between this repo and that other SOAP repo?
This is the SOAP Server Software. The other repository only has the communication templates between a Wii and WSC's server.

# Changelog
Versions on this software are based on goals. (e.g 0.2 works towards SQL support. 0.3 works towards NUS support, etc.)
## 0.2.x Kawauso
### 0.2.6
*This version of WiiSOAP Server was brought to you by Apfel. Thank you for your contribution.*
- Fixed error handling.
- Moved configuration example.
- Added `go.mod` for an easier installation.
- Changed `SQLPort` to `SQLAddress` in the `config.xml` file.
### 0.2.5
*This version of WiiSOAP Server was brought to you by Apfel. Thank you for your contribution.*
- Fixed lint errors.
- Uses Fprintf properly now.
- Uses `if err = action(); os.IsExist(err) {}` now. This makes error checks a little bit shorter.
- Changed `Port` to `Address` in the `config.xml` file.
### 0.2.4
- Added SQL skeleton.
- Edited config template.
### 0.2.3
- Added TODO Items.
- Improved Error Handling.
- Bug fixes.
> Fixed an issue where converting string to byte in the switch cases would cause the program to not compile. Since converting to byte in switch is not possible. The program should now compile without any errors.
- Organized the script to make it easier to read.
- No changes to SQL have been made in this update.
### 0.2.2
- Switched from if else to switch case. This makes the script cleaner, and makes the program faster.
- No changes to SQL have been made in this update.
### 0.2.1
- Added working Config.
- SQL now works. (In terms of opening a connection.)
- You can now choose what port to run WiiSOAP on.
### 0.2.0
- Added SQL Driver.
- SOAP-GO-OSC is now just called WiiSOAP.

## 0.1 (No Codename)
### 0.1.2
- Optimised the software. (Structures are now in their own file.)
### 0.1.1
- Added IAS Support.
### 0.1.0
- Added ECS Support.
