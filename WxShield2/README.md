# WxShield2

WxShield2 is a simple sketch that reads values off Sparkfun's 
[Weather Shield](https://www.sparkfun.com/products/12081) and emits JSON for the measured values.

# Exciting?!?!?

Ive known Arduino uses gcc/g++ as the compiler suite, and wanted to see how far I could streatch the compiler. The only real interesting thing is in helpers.h which uses function pointers to emit data.  One thing I ran into, which is probably a compiler toolchain limitation, is that any function pointer definitions must be added as a secondary .h header file in a project - you cannot do typedef's to function pointers in `*.ino` files and it work properly.

