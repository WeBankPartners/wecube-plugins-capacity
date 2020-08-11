<template>
  <div class="parameter-configuration">
    <Row>
      <Col span="3">
        <span class="param-title">{{$t('endpointAndMetric')}}</span>
      </Col>
      <Col span="21">
        <Form :label-width="20">
          <FormItem class="param-inline">
            <Select
              style="width:320px;"
              :placeholder="$t('placeholder.endpointSearch')"
              v-model="endpoint"
              filterable
              remote
              clearable
              @on-clear="getEndpointList('.')"
              :remote-method="getEndpointList"
              >
              <Option v-for="(option, index) in endpointList" :value="option.option_value" :key="index">{{option.option_text}}</Option>
              <Option value="moreTips" disabled>{{$t('tips.requestMoreData')}}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
            <Select 
              v-model="metric" 
              @on-open-change="getMetric" 
              style="width:320px" 
              :placeholder="$t('placeholder.metric')" 
              filterable
              clearable
              @on-clear="clearMetric">
                <Option v-for="item in metricList" :value="item.metric" :key="item.metric">{{ item.metric }}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
            <button @click="addParams()" :disabled="!endpoint || !metric" type="button" class="btn btn-confirm-skeleton-f">{{$t('add')}}</button>
          </FormItem>
          <div v-if="endpointWithMetric.length" class="params-display">
            <div v-for="(tag, tagIndex) in endpointWithMetric" :key="tagIndex" class="params-display-single">
              <Tag type="border" @on-close="removeTag(tagIndex)" closable color="primary">{{tag.endpoint}}:{{tag.metric}}</Tag>
            </div>
          </div>
        </Form>
      </Col>
    </Row>
    <Row>
      <Col span="3">
        <span class="param-title">{{$t('timeInterval')}}</span>
      </Col>
      <Col span="21">
      <Form :label-width="20">
        <FormItem class="param-inline">
          <DatePicker 
            type="datetimerange" 
            :value="dateRange" 
            @on-change="datePick"
            format="yyyy-MM-dd HH:mm:ss" 
            placement="bottom-start" 
            :placeholder="$t('placeholder.datePicker')" 
            style="width: 320px">
          </DatePicker>
        </FormItem>
        <FormItem class="param-inline">
          <button :disabled="endpointWithMetric.length && !dateRange[0]" @click="getChart" type="button" class="btn btn-confirm-f">{{$t('queryView')}}</button>
        </FormItem>
      </Form>
      </Col>
    </Row>
    <div :id="elId" class="echart" style="height:500px;width:1000px;box-shadow: 0 2px 20px 0 rgba(0,0,0,.11);margin-top:20px;margin-bottom:80px;"></div>
  </div>
</template>

<script>
import {generateUuid} from '@/assets/js/utils'
import {readyToDraw} from  '@/assets/config/chart-rely'
export default {
  name: '',
  data() {
    return {
      endpoint: '',
      endpointObject: null,
      endpointList: [],

      metric: '',
      metricList: [],

      endpointWithMetric: [],
      dateRange: ['', ''],

      chartData: null,
      elId: '',

      params: null
    }
  },
  created (){
    generateUuid().then((elId)=>{
      this.elId =  `id_${elId}`
    })
  },
  watch: {
    endpoint: function(val) {
      this.endpointObject = this.endpointList.find(ed => ed.option_value === val)
    }
  },
  mounted () {
    this.getEndpointList('.')
  },
  methods: {
    getChart () {
      if (this.dateRange[0] === this.dateRange[1]) {
        this.dateRange[1] = this.dateRange[1].replace('00:00:00', '23:59:59')
      }
      const start = Date.parse(this.dateRange[0])/1000 + ''
      const end = Date.parse(this.dateRange[1])/1000 + ''
      let params = []
      this.endpointWithMetric.forEach(single => {
        params.push({
          ...single,
          start,
          end
        })
      })
      this.$parent.params = params
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.getChart, params, (responseData) => {
        this.chartData = responseData.data
        const chartConfig = {eye: false,clear:true, zoomCallback: true}
        readyToDraw(this,responseData.data, 1, chartConfig)
        let xyAxis = []
        responseData.data.legend.forEach(_ => {
          xyAxis.push({
            label: _,
            value: _,
          })
        })
        this.$parent.xyAxis = xyAxis
      })
    },
    addParams () {
      if (this.endpointObject && this.metric) {
        this.endpointWithMetric.push({
          endpoint: this.endpointObject.option_value,
          metric: this.metric
        })
        this.endpointObject = null
        this.metric = ''
        this.endpoint = ''
        this.endpointList = []
        this.getEndpointList('.')
        this.$Message.success('OK')
      } else {
        this.$Message.error('error')
      }
    },
    removeTag (index) {
      this.endpointWithMetric.splice(index, 1)
    },
    getEndpointList(query) {
      this.metric = ''
      this.metricList = []
      let params = {
        search: query,
        search_type: 'endpoint'
      }
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('GET', this.$root.capacityApiCenter.getEndpoint, params, (responseData) => {
        this.endpointList = responseData.data
      })
    },
    clearMetric () {
      this.metric = ''
    },
    getMetric (val) {
      if (!val) return
      if (!this.endpointObject) {
        this.$Message.error(this.$t('tips.selectData'))
        return
      }
      let params = {
        type: this.endpointObject.type,
        search_type: 'metric'
      }
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('GET', this.$root.capacityApiCenter.getEndpoint, params, (responseData) => {
        this.metricList = responseData.data
      })
    },
    datePick (data) {
      this.dateRange = data
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.parameter-configuration {
  padding: 32px 40px;
}
.margin-left {
  margin-left: 20px;
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
.params-display {
  border: 1px solid #dbdbdb;
  margin-bottom: 24px;
  margin-left: 20px;
  width: 74%;
}
.params-display-single {
  padding: 4px;
}
</style>
