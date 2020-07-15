# keyboardmacro

## Introduction
Associate programs to keys of a secondary keyboard on linux.

## Permissions
Add your user to the input group.  
To know the correct group for your distro, you can run: 

```
ls -al /dev/input
```

Is is also possible to run it as root, but every command executed by the program will also be root.

## Configuration
The configuration is provided as a yaml file.  
```
devices:
    - device:
        name: 'usb-_USB_Keyboard-event-kbd' # id of the input device
        grab: true # Prevent other applications from handling the event
        bindings:
            30: # scan code for the key 'A'
                exec: 'firefox'
            31: # scan code for the key 'S'
                exec: 'notify-send "special message"'
    - device
        name: 'usb-_another_USB_Keyboard_'
        grab: false
        bindings:
            30:
                exec: 'email-app'
```

## Usage
Use the **-config** flag to provide the path to the configuration file.
```
./keyboardmacro -config my_config.yaml
```


The **-scan** flag will output every key press from every available input source. The output will list the device name and the scan code.  

```
# Example output
usb-_USB_Keyboard-event-kbd: 36
```