<template>
  <div class="display-result">
    <Row style="margin-bottom:16px">
      <Col span="3">
        <span class="param-title">Change Level</span>
      </Col>
      <Col span="21">
        <RadioGroup v-model="diyLevel" type="button">
          <template v-for="item in ['1','2','3']">
            <Radio :label="item" :style="{color: ['#19be6b', '#f90', '#ed4014'][item]}" :key="item">{{$t('level'+item)}}</Radio>
          </template>
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
        <div style="height: 200px;overflow: auto;padding:8px;background: #586b73;color:white">
          <div v-html="result.output"></div>
        </div>
      </Col>
    </Row>
    <Row style="margin-bottom:16px">
      <Col span="3">
        <span class="param-title">{{$t('formula')}}</span>
      </Col>
      <Col span="21">
        <Badge :text="$t('level'+result.level)" :type="['normal', 'success', 'warning', 'error'][result.level] || 'normal'">
          <div style="padding-right:16px">{{result.func_expr}}</div>
        </Badge>
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
      v-model="collectModel"
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
require('echarts/lib/component/toolbox');
export default {
  name: '',
  data() {
    return {
      isImportData: false,
      diyLevel: '0',
      result: {
        level: '0'
      },
      imgTips: ['photoTipsOne', 'photoTipsTwo', 'photoTipsThree', 'photoTipsFour'],
      collectModel: false,
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
  watch: {
    diyLevel: function (val) {
      this.diyLevel = val
      this.getRAnalyze()
    }
  },
  activated () {
    this.$parent.cachedCom.push(this.$vnode)
    this.isImportData = this.$parent.isImportData
    this.formulaParams = this.$parent.formulaParams
    this.getRAnalyze()
  },
  methods: {
    saveFormula () {
      this.collectModel = true
    },
    save () {
      this.result.level = Number(this.result.level)
      let params = {
        ...this.result,
        name: this.formInline.name,
        y_real: this.result.chart.data_series[0].data,
        y_func: this.result.chart.data_series[1].data,
      }
       this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.saveRAnalyze, params, () => {
        this.$Message.success('Success!')
      })
    },
    getRAnalyze () {
      let params = {
        min_level: Number(this.diyLevel)
      }
      if (this.isImportData) {
        params.excel = this.formulaParams
        params.guid = this.formulaParams.guid
        params.excel.enable = true
      } else {
        params.monitor = this.formulaParams
      }
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.getRAnalyze, params, (responseData) => {
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
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: 'none'
            }
          }
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
