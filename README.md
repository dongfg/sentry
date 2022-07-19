sentry
---
> keep your git repository up-to-date with webhook

### Development

```shell
git flow init

# feature
git flow feature start xxx
git flow feature finish -S

# release
git flow release start yyy
git flow release finish --nodevelopmerge

# after release
git checkout develop
git merge master
```