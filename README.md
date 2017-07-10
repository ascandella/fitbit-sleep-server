[![Build Status](https://travis-ci.org/sectioneight/fitbit-sleep-server.svg?branch=master)](https://travis-ci.org/sectioneight/fitbit-sleep-server)

Fitbit Sleep Server
===================

This is the code that powers [aidens.life](https://aidens.life/)

Systemd/Docker service
----------------------

I run my site in a Docker container under systemd. If you'd like to do something
similar, copy [my example service](systemd/example.service) to e.g.
`/etc/systemd/system` and then run `systemctl daemon-reload`.

Deploying
---------

I use [Caddy](https://caddyserver.com) and you probably should too. This is how
I have my host configured to auto-deploy via `systemctl`, `make`, and Webhooks:

```plain
aidens.life {
  root /opt/ai-life

  git {
    repo https://github.com/sectioneight/fitbit-sleep-server.git
    hook /webhook <redacted>
    then sudo /bin/systemctl restart ai-life
  }

  proxy / localhost:3030 {
    except /webhook
    transparent
  }
}
```

This requires an appropriate entry in `sudoers.d`, such as:

```bash
caddy ALL=NOPASSWD: /bin/systemctl restart ai-life
```

Where `caddy` is the user your caddy server is running under.
