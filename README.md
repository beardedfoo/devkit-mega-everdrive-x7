# devkit-mega-everdrive-x7
Development utilities for the Sega Genesis/Megadrive with Mega Everdrive X7 by @krikzz

The Sega Megadrive (known as Genesis in US) is a video game console released by Sega in 1988. Modern development cartridges have been developed for the system, and one such cartridge is the Mega Everdrive X7, which features a USB Serial TTY and an SD memory card. The Mega Everdrive X7 allows game roms to be loaded from a filesystem on the SD card, but also allows game data to be sent to the game console over the USB serial TTY. After game software loads the USB serial TTY can also be used by the Megadrive to communicate with host programs on a computer, enabling developers to write code for debugging as well as adding new features such as internet connectivity to the game console.

The software distributed for the Mega Everdrive X7 is windows-only and written in C#, but offers some compatibility with other systems via the mono project. This project is an attempt increase accessibilty of development using the Mega Everdrive X7 by creating natively cross-platform utilities.
