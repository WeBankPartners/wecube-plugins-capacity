FROM r-base:4.0.3

ENV BASE_HOME=/app/capacity

RUN mkdir -p $BASE_HOME $BASE_HOME/conf $BASE_HOME/logs $BASE_HOME/public

ADD build/start.sh $BASE_HOME/
ADD build/stop.sh $BASE_HOME/
ADD build/default.json $BASE_HOME/conf/
ADD server/server $BASE_HOME/
ADD server/public/index.html $BASE_HOME/public/
ADD server/public/favicon.ico $BASE_HOME/public/
ADD server/public/js $BASE_HOME/public/js
ADD server/public/css $BASE_HOME/public/css
ADD server/public/fonts $BASE_HOME/public/fonts
ADD server/public/img $BASE_HOME/public/img
ADD server/conf/template.r $BASE_HOME/conf/

WORKDIR $BASE_HOME
ENTRYPOINT ["/bin/sh", "start.sh"]