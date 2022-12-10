 #go build *.go &&  sudo -E setcap cap_net_raw=+ep ./metrics && ./metrics
 go build -o app ./... &&  sudo -E setcap cap_net_raw=+ep ./app && ./app
