# delete old png images from outputs folder
find /home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/outputs/*.png -daystart -mtime +1 -delete