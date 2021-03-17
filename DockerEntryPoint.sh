#!/bin/sh
set -e

mkdir -p tmp public/img/banner public/img/promotion public/img/reward public/img/luckydraw public/img/luckydrawReward public/highlight/header public/highlight/order public/img/promotion public/img/promotion/order public/img/promotion/brand public/img/pointrule/brand public/img/pointrule/product public/img/user
chmod 777 tmp public/img/banner public/img/promotion public/img/reward public/img/luckydraw public/img/luckydrawReward public/highlight/header public/highlight/order public/img/promotion public/img/promotion/order public/img/promotion/brand public/img/pointrule/brand public/img/pointrule/product public/img/user

cp -f /usr/share/zoneinfo/${TIMEZONE} /etc/localtime; \
echo "${TIMEZONE}" >  /etc/timezone;

exec "$@"