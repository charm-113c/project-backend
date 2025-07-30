# This Dockerfile is only used for development purposes
# DO NOT USE IN PROD
# This image is used to run dev mode in a standardised environment.

# Base image
FROM golang:1.23.4 AS builder

# Install Make to run Make commands
RUN apt-get update && apt-get install -y make

# Set working directory
WORKDIR /app

# Copy code in there
COPY . .

# To make the final image lightweight,
# we'll build the code in this "builder" image
# and only save the binary + config in the final image

RUN make build
# As per make build, the binary is saved in bin/

# Final image
FROM debian:bookworm

# Set to development mode
ENV DEV_MODE=true
ENV ENV_PATH="config.env"

# Open port
EXPOSE 7777

WORKDIR /app/bin

COPY --from=builder /app/bin/backend /app/config/config.env ./

# ENTRYPOINT ["/bin/bash"]
CMD ["./backend"]
