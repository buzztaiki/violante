Virus scan server and cilent using [VirusTotal](https://www.virustotal.com).

## Motivation

I hope that I want to rescan files that reported by ClamAV using VirusTotal, and then I want to notify to slack if it is really virus.


## Install

### Archlinux

```console
> git clone https://github.com/buzztaiki/violante
> cd violante/archlinux
> makepkg -si
> sudo systemctl enable violante-server
> sudo systemctl edit violante-server
[Service]
Environment=VT_API_KEY=your_virustotal_api_key
Environment=SLACK_WEBHOOK_URL=https://hooks.slack.com/services/your_slack_webhook
Environment=SLACK_CHANNEL=#your_slack_channel
> sudo systemctl start violante-server
```

### Another OS

```
> go get github.com/buzztaiki/violante/tree/master/cmd/violante
> go get github.com/buzztaiki/violante/tree/master/cmd/violante-server
```

and then?


## License
MIT
