# devkit-mega-everdrive-x7
Development utilities for the Sega Genesis/Megadrive with Mega Everdrive X7 by @krikzz

The Sega Megadrive (known as Genesis in US) is a video game console released by Sega in 1988. Modern development cartridges have been developed for the system, and one such cartridge is the Mega Everdrive X7, which features a USB Serial TTY and an SD memory card. The Mega Everdrive X7 allows game roms to be loaded from a filesystem on the SD card, but also allows game data to be sent to the game console over the USB serial TTY. After game software loads the USB serial TTY can also be used by the Megadrive to communicate with host programs on a computer, enabling developers to write code for debugging as well as adding new features such as internet connectivity to the game console.

The official software distributed for the Mega Everdrive X7 is windows-only and written in C#, but offers some compatibility with other systems via the mono project. This project is an attempt increase accessibilty of development using the Mega Everdrive X7 by creating natively cross-platform utilities.

I have managed to create a utility called megaedx7-run, written in golang using a cross-platform serial library, which can interact with the Sega Megadrive/Mega Everdrive X7. This program can load arbitrary game roms over USB and execute them on the game console in various modes as supported by the Mega Everdrive. These modes are os, cd, sms
* `sms`: Sega Master System ROM (untested)
* `md`: Megadrive/Genesis ROM
* `cd`: Mega CD / Sega CD ROM; may be problematic due to address space collisions between the Mega CD add-on and the Mega Everdrive X7 
* `m10`: Unknown
* `os`: Unknown; seems to load Mega Everdrive X7 firmware ROMs
* `ssf`: Unknown; seems to be indicated by the string 'SSF ' at ROM offset 0x105

The Mega Everdrive X7 cartridge seems to expect raw ROM dumps; other formats have not worked in my testing. ROMs with the ASCII string 'SEGA' at offset 0x100 have worked in my testing.

![megaedx7-run usage](images/megaedx7-help.png?raw=true)
![megaedx7-run example](images/megaedx7-run-cmd.png?raw=true)

