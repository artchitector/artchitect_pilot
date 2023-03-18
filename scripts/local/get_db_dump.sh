# go to memory server and download dump locally
mkdir -p ~/dumps/
mkdir -p ~/dumps/db/
# get last archive from memory server
echo "Start get dump from memory server"
lastDump=$(ssh memory "cd ~/dumps/db && ls -tr | tail -1")
echo "Last dump is $lastDump on memory server"
path=~/dumps/db/$lastDump
if [ -e "$path" ]; then
  echo "dump $lastDump already downloaded. exit"
else
  scp memory:~/dumps/db/$lastDump $path
  echo "successfully downloaded dump $lastDump to $path"
fi

