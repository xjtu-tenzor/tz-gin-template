local_version=$(cat version.txt)
remote_version=$(curl -s https://raw.githubusercontent.com/Asice-Cloud/tz-gin-template/main/version.txt)
if [ "$local_version" = "$remote_version" ]; then
  echo "Already latest version: $local_version"
else
  echo "New version available: $remote_version (local: $local_version)"
fi