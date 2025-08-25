# Use Debian x86_64 as base (via multiarch/qemu)
FROM --platform=linux/amd64 debian:bullseye

# Install assembler + compiler toolchain and dependencies
RUN apt-get update && apt-get install -y \
    nasm \
    binutils \
    make \
    gdb \
    build-essential \
    curl \
    wget \
    tar \
    && rm -rf /var/lib/apt/lists/*

# Set working directory inside container
WORKDIR /app

# Default command: interactive shell
CMD ["/bin/bash"]
