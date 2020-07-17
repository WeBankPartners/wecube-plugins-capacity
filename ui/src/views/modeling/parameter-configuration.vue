<template>
  <div class="parameter-configuration">
    <Row>
        <Col span="3">
          <span class="param-title">指标配置</span>
        </Col>
        <Col span="21">
        <Form :label-width="60">
          <FormItem class="param-inline" label="对象">
            <Select
              style="width:270px;"
              placeholder="请选择对象"
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
          <FormItem class="param-inline" label="指标">
            <Select v-model="metric" @on-open-change="getMetric" style="width:270px" placeholder="请选择指标" filterable>
              <Option v-for="item in metricList" :value="item.metric" :key="item.metric">{{ item.metric }}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
              <button class="btn btn-confirm-f">增加</button>
            </FormItem>
        </Form>
        </Col>
    </Row>
    <!-- <Row>
      <Col span="3">
        <span class="param-title">时间区间</span>
      </Col>
      <Col span="21">
      <Form :label-width="60">
        <FormItem class="param-inline" label="对象">
          <Select
            style="width:270px;"
            placeholder="请选择对象"
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
        <FormItem class="param-inline" label="指标">
          <Select v-model="metric" @on-open-change="getMetric" style="width:270px" placeholder="请选择指标" filterable>
            <Option v-for="item in metricList" :value="item.metric" :key="item.metric">{{ item.metric }}</Option>
          </Select>
        </FormItem>
        <FormItem class="param-inline">
            <button class="btn btn-confirm-f">增加</button>
          </FormItem>
      </Form>
      </Col>
    </Row> -->
  </div>
</template>

<script>
export default {
  name: '',
  data() {
    return {
      endpoint: '',
      endpointObject: null,
      endpointList: [],

      metric: '',
      metricList: []
    }
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
    getEndpointList(query) {
      let params = {
        search: query,
        search_type: 'endpoint'
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.getEndpoint, params, (responseData) => {
        this.endpointList = responseData.data
      })
    },
    getMetric (val) {
      if (!val) return
      if (!this.endpointObject) {
        this.$Message.error('请先选择对象！')
        return
      }
      let params = {
        type: this.endpointObject.type,
        search_type: 'metric'
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.getEndpoint, params, (responseData) => {
        this.metricList = responseData.data
      })
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.parameter-configuration {
  padding: 32px 40px;
}
.ivu-form-item {
  margin-bottom: 10px;
}
.param-title {
  line-height: 40px;
}
.param-inline {
  display: inline-block;
}
</style>
