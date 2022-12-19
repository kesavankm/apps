 #go build *.go &&  sudo -E setcap cap_net_raw=+ep ./metrics && ./metrics
 go build -o ./bin/app ./src &&  sudo -E setcap cap_net_raw=+ep ./bin/app && ./bin/app
