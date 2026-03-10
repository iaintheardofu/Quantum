# SomaOS: MABEL x8C C-Agent Makefile
# Target: Zynq-7000 (ARM Cortex-A9)

CC = arm-linux-gnueabihf-gcc
CFLAGS = -O2 -Wall
LDFLAGS = 

TARGET = soma_agent
SRC = src/soma_agent.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(SRC) -o $(TARGET) $(LDFLAGS)

clean:
	rm -f $(TARGET)

.PHONY: all clean
