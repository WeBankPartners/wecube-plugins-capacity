<template>
  <div class="modeling-container">
    <div class="modeling-steps">
      <Row style="text-align: center">
        <Col span="3">
          <span class="step-title">{{$t('menu.capacityForecast')}}</span>
        </Col>
      </Row>
    </div>
    <div class="operation">
      <div class="operation-zone">
        <Row>
          <Col span="3">
            <span class="param-title">{{$t('favoritesList')}}</span>
          </Col>
          <Col span="21">
            <Form>
              <FormItem class="param-inline">
                <Select v-model="favorite" style="width:320px" :placeholder="$t('placeholder.metric')" filterable>
                  <Option v-for="item in favoritesList" :value="item.guid" :key="item.guid">{{ item.name }}</Option>
                </Select>
              </FormItem>
              <FormItem class="param-inline">
                <button @click="getFavoriteDetail()" :disabled="!favorite" type="button" class="btn btn-confirm-skeleton-f">{{$t('search')}}</button>
              </FormItem>
            </Form>
          </Col>
        </Row>

        <div>
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
            </Col>
          </Row>
          <Row style="margin-bottom:16px">
            <Col span="3">
              <span class="param-title">Output</span>
            </Col>
            <Col span="21">
              <div  style="height: 200px;overflow: auto;background: skyblue;padding:8px">
                <div v-html="result.output"></div>
              </div>
            </Col>
          </Row>
          <Row style="margin-bottom:16px">
            <Col span="3">
              <span class="param-title">{{$t('formula')}}</span>
            </Col>
            <Col span="21">
              {{result.func_expr}}
            </Col>
          </Row>
          <div style="text-align: center">
            <template v-for="(img) in result.images">
              <img :src="'http://129.204.99.160:19696/capacity/'+ img" :key="img" alt="">
            </template>
          </div>
          <div>
            <div id="graph" class="echart" style="height:500px;width:1000px;box-shadow: 0 2px 20px 0 rgba(0,0,0,.11);margin-top:40px"></div>
          </div>
          <div v-if="tableData.length" style="margin: 16px 20px 16px 0">
            <Table 
              border 
              ref="selection" 
              :columns="columns" 
              :data="tableData">
            </Table>
          </div>
          <Row v-if="result.level" style="margin-top:16px">
            <Col span="3">
              <span class="param-title">输入预测值</span>
            </Col>
            <Col span="21">
              <div v-for="(item, intemIndex) in inputArray" :key="intemIndex" class="user-array">
                <template v-for="(key, keyIndex) in Object.keys(item)">
                  <InputNumber :min="1" :step="0.1" v-model="item[key]" :key="key + intemIndex + keyIndex" class="user-input"></InputNumber>
                </template>
                <Icon type="md-trash" @click="deleteRow(intemIndex)" v-if="inputArray.length != 1" class="operation-icon-delete" />
                <Icon type="ios-add" @click="addRow" v-if="intemIndex + 1 === inputArray.length" class="operation-icon-add" />
              </div>
              <button @click="prediction" type="button" style="margin-top:8px" class="btn btn-confirm-skeleton-f">预测</button>
            </Col>
          </Row>
        </div>
      </div>
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
      favorite: null,
      favoritesList: [],
      selectedData: null,

      showResult:false,
      result: {
        level: ''
      },

      inputTmp: null,
      inputArray: [],

      columns: [],
      tableData: []
    }
  },
  mounted () {
    this.getFavoritesList()
  },
  methods: {
    getFavoritesList () {
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.getFavoritesList, '', (responseData) => {
        this.favoritesList = responseData.data
      })
    },
    getFavoriteDetail () {
      this.inputTmp = null
      this.inputArray = []
      const params = {
        guid: this.favorite
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.favoriteDetail, params, (responseData) => {
        this.result = responseData.data
        this.result.level = responseData.data.level + ''
        this.drawChart(responseData.data.chart)
        const len = responseData.data.func_x.length
        let tmp = {}
        for (let i = 0; i < len; i++) {
          tmp['key' + i] = null
        }
        this.inputTmp = JSON.parse(JSON.stringify(tmp))
        this.inputArray.push(tmp)
        this.showResult = true
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
    },
    addRow () {
      this.inputArray.push(JSON.parse(JSON.stringify(this.inputTmp)))
    },
    deleteRow (index) {
      this.inputArray.splice(index, 1)
    },
    prediction () {
      let array_ = []
      
      for (let item of this.inputArray) {
        const tmp = Object.values(item).some(i => {
          return i === null
        })
        if (tmp) {
          this.$Message.warning('预测参数不允许为空！')
          return
        }
        array_.push(Object.values(item))
      }
      let params = {
        guid: this.favorite,
        add_data: array_
      }
      this.columns = []
      this.tableData = []
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.favoriteCalc, params, (responseData) => {
        this.drawChart(responseData.data.chart)
        responseData.data.table.title.forEach(i => {
          this.columns.push({
            title: i,
            key: i
          })
        })
        this.tableData = responseData.data.table.data
      })
    }
  },
  components: {
  },
}
</script>

<style scoped lang="less">
.modeling-container {
  font-size: 16px;
}
.modeling-steps {
  background: #f1f1f1;
  width: 1180px;
  margin: 0 auto;
  padding: 16px 0;
  background: white;
}
.step-icon {
  padding: 4px;
  border-radius: 50%;
  border: 1px solid #3ba3ff;
  color: #3ba3ff;
}
.step-title {
  padding: 8px;
  font-weight: 400;
}
.operation {
  background: #f0f2f5;
  padding-top: 20px;
}
.operation-zone {
  width: 1100px;
  min-height: calc(~"100vh - 140px");
  margin: 0 auto;
  background: #ffffff;
  padding: 32px 40px;
}
.step-control {
  position: fixed;
  bottom: 0;
  height: 80px;
  width: 1100px;
  padding: 24px 68px;
  border-top: 1px solid #dbe3e4;
  box-shadow: 0 -4px 4px -2px #e4e9f0;
  background: white;
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
// .user-array {
//   margin-bottom: 8px;
// }
.user-input {
  margin-right: 8px;
}
.operation-icon-delete {
  font-size: 18px;
  border: 1px solid #ed4014;
  color: #ed4014;
  border-radius: 4px;
  width: 30px;
  line-height: 30px;
  margin: 6px;
  vertical-align: middle;
}
.operation-icon-add {
  font-size: 24px;
  border: 1px solid #0080FF;
  color: #0080FF;
  border-radius: 4px;
  width: 30px;
  line-height: 30px;
  margin: 6px;
  vertical-align: middle;
}
</style>
