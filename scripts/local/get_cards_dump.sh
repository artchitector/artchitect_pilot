# go to memory server and download dump locally
mkdir -p ~/dumps/
mkdir -p ~/dumps/cards/
# get last cards zip file from memory server
echo "start get cards dump from memory server"
# find last dump older than 5 minutes (zipped fully)
lastDump=$(ssh memory "find ~/dumps/cards -type f \( -mmin +10 \) -printf \"%T@ %f\n\" | sort | cut -d' ' -f2 | tail -1")
echo "last cards dump is $lastDump on memory server"

path=~/dumps/cards/$lastDump
if [ -e "$path" ]; then
  echo "cards $lastDump already downloaded. exit"
else
  fullPath="/root/dumps/cards/$lastDump"
  scp memory:$fullPath $path
  echo "successfully downloaded cards zip $lastDump to $path"
fi