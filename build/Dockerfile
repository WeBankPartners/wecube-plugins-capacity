FROM r-base

ENV BASE_HOME=/app/capacity

RUN mkdir -p $BASE_HOME $BASE_HOME/conf $BASE_HOME/logs

ADD ./start.sh $BASE_HOME/
ADD ./stop.sh $BASE_HOME/
ADD ./default.json $BASE_HOME/conf/
ADD ../server/server $BASE_HOME/
ADD ../server/public $BASE_HOME/
ADD ../server/conf/template.r $BASE_HOME/conf/

WORKDIR $BASE_HOME
ENTRYPOINT ["/bin/sh", "start.sh"]