#!/bin/bash

while true; do
    echo "Hello from the custom service! at $(date)" >> /home/amit/personal/blog-content/mastering-linux-services/log_message.log
    sleep 10  # Wait for 10 seconds
done
