<template>
  <div class="modeling-container">
    <div style="background:white">
      <div class="modeling-steps">
        <Row style="text-align: center">
          <Col span="4">
            <span class="step-title">{{$t('capacityForecast')}}</span>
          </Col>
        </Row>
      </div>
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
            <div id="graph" class="echart" style="height:500px;width:1020px;box-shadow: 0 2px 20px 0 rgba(0,0,0,.11);margin-top:40px"></div>
          </div>
          <div v-if="tableData.length" style="margin: 16px 20px 16px 0">
            <Table 
              border 
              size="small"
              ref="selection" 
              :columns="columns" 
              :data="tableData">
            </Table>
          </div>
          <Row v-if="result.level" style="margin-top:16px">
            <Col span="3">
              <span class="param-title">{{$t('predictedValue')}}</span>
            </Col>
            <Col span="21">
               <Tag v-if="targetObject" color="primary">{{targetObject}}</Tag>
              <div v-for="(item, intemIndex) in inputArray" :key="intemIndex" class="user-array">
                <template v-for="(key, keyIndex) in Object.keys(item)">
                  <InputNumber 
                    :min="0" 
                    :step="0.1" 
                    v-model="item[key]" 
                    :placeholder="inputPlaceholder[keyIndex]"
                    :key="key + intemIndex + keyIndex" 
                    class="user-input"
                    >
                  </InputNumber>
                </template>
                <Icon type="md-trash" @click="deleteRow(intemIndex)" v-if="inputArray.length != 1" class="operation-icon-delete" />
                <Icon type="ios-add" @click="addRow" v-if="intemIndex + 1 === inputArray.length" class="operation-icon-add" />
              </div>
              <button @click="prediction" type="button" style="margin-top:8px" class="btn btn-confirm-skeleton-f">{{$t('prediction')}}</button>
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
require('echarts/lib/component/toolbox')
export default {
  name: '',
  data() {
    return {
      favorite: null,
      favoritesList: [],
      selectedData: null,

      imgTips: ['photoTipsOne', 'photoTipsTwo', 'photoTipsThree', 'photoTipsFour'],
      showResult:false,
      result: {
        level: '0'
      },

      inputTmp: null,
      inputArray: [], // 预测值
      inputPlaceholder: [], // 预测值提示信息
      targetObject: null, // 待预测机器

      columns: [],
      tableData: []
    }
  },
  mounted () {
    this.getFavoritesList()
  },
  methods: {
    getFavoritesList () {
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('GET', this.$root.capacityApiCenter.getFavoritesList, '', (responseData) => {
        this.favoritesList = responseData.data
      })
    },
    getFavoriteDetail () {
      this.inputTmp = null
      this.inputArray = []
      const params = {
        guid: this.favorite
      }
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('GET', this.$root.capacityApiCenter.favoriteDetail, params, (responseData) => {
        this.result = responseData.data
        this.result.level = responseData.data.level + ''
        this.drawChart(responseData.data.chart)
        const len = responseData.data.func_x.length
        this.inputPlaceholder = responseData.data.legend_x.map(item => {
          return item.split(':')[1]
        })
        this.targetObject = responseData.data.legend_x[0].split(':')[0]
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
          this.$Message.warning(this.$t('tips.forecastEmpty'))
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
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.favoriteCalc, params, (responseData) => {
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
  width: 1180px;
  margin: 0 auto;
  padding: 16px 0;
}
.step-icon {
  padding: 4px;
  border-radius: 50%;
  border: 1px solid #3ba3ff;
  color: #3ba3ff;
}
.step-title {
  font-size: 18px;
  padding: 8px;
  font-weight: 500;
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
  width: 270px;
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
