# delete old png images from outputs folder
find /home/artchitector/invoke-ai/invokeai_v2.3.0/invokeai/outputs/*.png -daystart -mtime +1 -delete