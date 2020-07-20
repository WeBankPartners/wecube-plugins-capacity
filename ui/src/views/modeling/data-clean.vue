<template>
  <div class="data-clean">
    <Row>
      <Col span="3">
        <span class="param-title">选择轴坐标</span>
      </Col>
      <Col span="21">
        <Form :label-width="20">
          <FormItem class="param-inline">
            <Select v-model="xValue" style="width:320px" :placeholder="$t('placeholder.xAxis')" filterable multiple>
              <Option v-for="item in xOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
            <Select v-model="yValue" style="width:320px" :placeholder="$t('placeholder.yAxis')" filterable>
              <Option v-for="item in yOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
            </Select>
          </FormItem>
        </Form>
      </Col>
    </Row>
    <Row>
      <Col span="21" offset="3">
        <button :disabled="!(xValue.length && yValue)" @click="getData" type="button" class="btn btn-confirm-f margin-left">查询数据</button>
      </Col>
    </Row>
  </div>
</template>

<script>
export default {
  name: '',
  data() {
    return {
      xValue: [],
      xOptions: [],

      yValue: null,
      yOptions: []
    }
  },
  props: ['params', 'xyAxis'],
  created () {
    this.xOptions = JSON.parse(JSON.stringify(this.xyAxis))
    this.yOptions = JSON.parse(JSON.stringify(this.xyAxis))
  },
  methods: {
    getData () {
      let params = {
        config: this.params,
        legend_x: this.xValue,
        legend_y: this.yValue,
        aggregate: 'none'
      }
      this.$root.$httpRequestEntrance.httpRequestEntrance('POST', this.$root.apiCenter.getRData, params, (responseData) => {
        console.log(responseData)
      })
    }
  },
  components: {},
}
</script>

<style scoped lang="less">
.data-clean {
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
</style>
