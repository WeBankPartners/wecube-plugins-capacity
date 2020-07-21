<template>
  <div class="display-result">
    <Row>
      <Col span="3">
        <span class="param-title">Level</span>
      </Col>
      <Col span="21">
      </Col>
    </Row>
    <!-- <div>level: {{result.level}}</div> -->
    <!-- <div>output: 
      <div v-html="result.output"></div>
    </div>
    <div>func_expr: {{result.func_expr}}</div> -->
    <!-- <div>
      <img src="http://129.204.99.160:19696/capacity/r_images/R_1595313716221396518/rp001.png" alt="">
      <img src="http://129.204.99.160:19696/capacity/r_images/R_1595313716221396518/rp002.png" alt="">
      <img src="http://129.204.99.160:19696/capacity/r_images/R_1595313716221396518/rp003.png" alt="">
      <img src="http://129.204.99.160:19696/capacity/r_images/R_1595313716221396518/rp004.png" alt="">
    </div> -->
    <div>
      <div id="graph" class="echart" style="height:500px;width:1000px;box-shadow: 0 2px 20px 0 rgba(0,0,0,.11);margin-top:40px"></div>
    </div>
  </div>
</template>

<script>
require('echarts/lib/chart/line');
const echarts = require('echarts/lib/echarts')
export default {
  name: '',
  data() {
    return {
      formulaParams: {
        "config":[
          {
            "end":"1594828800",
            "endpoint":"VM_0_16_centos_192.168.0.16_host",
            "metric":"mem.used.percent",
            "start":"1594742400",
            "agg":"p95"
          }
        ],
        "legend_x":[
          "VM_0_16_centos_192.168.0.16_host:mem.used.percent"
        ],
        "legend_y":"VM_0_16_centos_192.168.0.16_host:mem.used.percent",
        "remove_list":[]
      },
      result: {
        level: null
      }
    }
  },
  // props: ['formulaParams'],
  mounted () {
    this.getRAnalyze()
  },
  methods: {
    getRAnalyze () {
      let params = {
        monitor: this.formulaParams
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.getRAnalyze, params, (responseData) => {
        this.result = responseData.data
        this.drawChart(responseData.data.chart)
      })
    },
    drawChart (config) {
      console.log(config)
      let myChart = echarts.init(document.getElementById('graph'))
      // let option = {
      //   xAxis: {
      //     data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
      //   },
      //   yAxis: {
      //   },
      //   series: [{
      //     data: [820, 932, 901, 934, 1290, 1330, 1320],
      //     type: 'line',
      //     smooth: true
      //   },
      //   {
      //     data: [820, 32, 901, 34, 90, 1330, 320],
      //     type: 'line',
      //     smooth: true
      //   }]
      // }
      let option = {
        legend: {
          bottom: 5,
          left: 'center',
          data: config.legend
        },
        xAxis: config.xaxis,
        yAxis: {
        },
        series: config.data_series
      }
      console.log(option)
      myChart.setOption(option)
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.display-result {
  padding: 32px 40px;
}
.ivu-form-item {
  margin-bottom: 24px;
}
.param-title {
  line-height: 32px;
}
.param-inline {
  display: inline-block;
}
.margin-left {
  margin-left: 20px;
}
</style>
