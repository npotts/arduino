

PROJECT_DIR   := $(shell pwd)
ARDUINO_DIR   = /usr/share/arduino
ARDMK_DIR     = /usr/share/arduino
AVR_TOOLS_DIR = /usr
USER_LIB_PATH :=  $(PROJECT_DIR)/../vendor

BOARD_TAG         = uno
MONITOR_BAUDRATE  = 115200
AVR_TOOLS_DIR     = /usr

AVRDDUDE          = /usr/bin/avrdude    # Path to avrdude directory.
CPPFLAGS          = -pedantic -Wall -Wextra
ARDUINO_LIBS      = Wire OneWire ArduinoJson MPL3115A2/src HTU21D/src SimpleDHT

MONITOR_PORT      = /dev/tty.USB*

### don't touch this
CURRENT_DIR       = $(shell basename $(CURDIR))
### OBJDIR
### This is were you put the binaries you just compile using 'make'
OBJDIR            = $(PROJECT_DIR)/bin/$(BOARD_TAG)/$(CURRENT_DIR)

### path to Arduino.mk, inside the ARDMK_DIR, don't touch.
include $(ARDMK_DIR)/Arduino.mk

