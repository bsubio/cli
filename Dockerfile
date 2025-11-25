FROM ubuntu:22.04

# Install ca-certificates for HTTPS
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the build
COPY bsubio /usr/local/bin/bsubio

# Set bsubio as the entrypoint
ENTRYPOINT ["/usr/local/bin/bsubio"]

# Default to help command
CMD ["help"]
