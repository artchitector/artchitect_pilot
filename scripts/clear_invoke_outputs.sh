# delete old png images from outputs folder
find /home/artchitector/invokeai/outputs/*.png -daystart -mtime +1 -delete