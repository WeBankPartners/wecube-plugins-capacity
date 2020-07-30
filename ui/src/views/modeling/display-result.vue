<template>
  <div class="display-result">
    <Row style="margin-bottom:16px">
      <Col span="3">
        <span class="param-title">Level</span>
      </Col>
      <Col span="21">
        <RadioGroup v-model="result.level" type="button">
          <Radio label="0" :disabled="result.level!='0'">{{$t('level0')}}</Radio>
          <Radio label="1" :disabled="result.level!='1'">{{$t('level1')}}</Radio>
          <Radio label="2" :disabled="result.level!='2'">{{$t('level2')}}</Radio>
          <Radio label="3" :disabled="result.level!='3'">{{$t('level3')}}</Radio>
        </RadioGroup>
        <span style="float:right">
          <button type="button" class="btn btn-confirm-f" @click="saveFormula" :disabled="result.level === '0'">{{$t('save')}}</button>
        </span>
      </Col>
    </Row>
    <Row style="margin-bottom:16px">
      <Col span="3">
        <span class="param-title">Output</span>
      </Col>
      <Col span="21">
        <div  style="height: 200px;overflow: auto;padding:8px;background: #586b73;color:white">
          <div v-html="result.output"></div>
        </div>
      </Col>
    </Row>
    <Row style="margin-bottom:16px">
      <Col span="3">
        <span class="param-title">公式</span>
      </Col>
      <Col span="21">
        {{result.func_expr}}
      </Col>
    </Row>
    <div style="text-align: center">
      <template v-for="(img, index) in result.images">
        <Tooltip placement="top" :key="index">
          <img :src="'/capacity/'+ img" :key="img" alt="">
          <div slot="content" style="white-space:normal">
            {{$t(imgTips[index])}}
          </div>
        </Tooltip>
      </template>
    </div>
    <div>
      <div id="graph" class="echart" style="height:500px;width:1000px;box-shadow: 0 2px 20px 0 rgba(0,0,0,.11);margin-top:40px"></div>
    </div>
    <Modal
      v-model="modal1"
      :title="$t('favorite')"
      @on-ok="save">
      <Form :model="formInline" :rules="ruleInline" :label-width="80">
        <FormItem :label="$t('name')" prop="name">
          <Input v-model="formInline.name" placeholder="Enter something..."></Input>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<script>
require('echarts/lib/chart/line');
const echarts = require('echarts/lib/echarts')
export default {
  name: '',
  data() {
    return {
      result: {
        level: '0'
      },

      imgTips: ['photoTipsOne', 'photoTipsTwo', 'photoTipsThree', 'photoTipsFour'],
      modal1: false,
      formInline: {
        name: ''
      },
      ruleInline: {
        name: [
          { required: true, message: 'Please fill in the name', trigger: 'blur' }
        ]
      }
    }
  },
  activated () {
    this.formulaParams = this.$parent.formulaParams
    this.getRAnalyze()
  },
  methods: {
    saveFormula () {
      this.modal1 = true
    },
    save () {
      this.result.level = Number(this.result.level)
      let params = {
        ...this.result,
        name: this.formInline.name,
        y_real: this.result.chart.data_series[0].data,
        y_func: this.result.chart.data_series[1].data,
      }
       this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.saveRAnalyze, params, () => {
        this.$Message.success('Success!')
      })
    },
    getRAnalyze () {
      let params = {
        monitor: this.formulaParams
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.getRAnalyze, params, (responseData) => {
        this.result = responseData.data
        this.result.level = responseData.data.level + ''
        this.drawChart(responseData.data.chart)
      })
    },
    drawChart (config) {
      let myChart = echarts.init(document.getElementById('graph'))
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
      myChart.setOption(option)
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.display-result {
  padding: 32px 40px;
  padding-bottom: 120px;
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
