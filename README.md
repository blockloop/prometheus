# Prometheus 

This is a fork of Prometheus that enables remote_write. It is used by @digitalocean/insights.

The master branch is kept up to date with the latest stable tag of upstream prometheus. All fork changes are written to
the `remote_write` branch. 

To make changes or fixes to remote_write capability, push a PR to the remote_write branch.

To update the master branch from upstream make sure you have the upstream remote added locally

```
git remote add upstream git@github.com:prometheus/prometheus.git
```

Once you have the remote you can pull from upstream, checkout the latest stable tag, and then rebase it into
the remote_write branch

```
git checkout master
git pull upstream master
git reset --hard <latest stable tag, i.e. v2.8.1>
git checkout remote_write
git rebase master
```

Once you have rebased master, build it with docker

```
docker build . -t digitalocean/prometheus:<latest stable tag, i.e. v2.8.1>
docker push . -t digitalocean/prometheus:<tag>
```
