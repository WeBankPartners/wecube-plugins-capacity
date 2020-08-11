<template>
  <div class="data-clean">
    <Row>
      <Col span="3">
        <span class="param-title">{{$t('axisCoordinates')}}</span>
      </Col>
      <Col span="21">
        <Form :label-width="20">
          <FormItem class="param-inline">
            <Select v-model="xValue" style="width:320px" :placeholder="$t('placeholder.xAxis')" filterable multiple>
              <Option v-for="item in xOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
            <Select 
              v-model="yValue" 
              style="width:320px" 
              :placeholder="$t('placeholder.yAxis')" 
              filterable
              clearable
              @on-clear="clearYValue">
              <Option v-for="item in yOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
            </Select>
          </FormItem>
          <FormItem class="param-inline">
            <Select v-model="aggregateValue" style="width:90px" :placeholder="$t('placeholder.aggregate')" filterable>
              <Option v-for="item in aggregateOption" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
          </FormItem>
        </Form>
      </Col>
    </Row>
    <Row>
      <Col span="21" offset="3">
        <button :disabled="!(xValue.length && yValue)" @click="getData" type="button" class="btn btn-confirm-f margin-left">{{$t('searchData')}}</button>
      </Col>
    </Row>
    <div style="height:500px;margin-top:20px;margin-bottom:125px">
      <Table 
        border 
        size="small"
        ref="selection" 
        :columns="columns" 
        @on-select="selectData"
        @on-select-cancel="selectCancle"
        @on-select-all="selectAllData"
        @on-select-all-cancel="selectAllCancle"
        :data="data"></Table>
      <div style="float: right;padding: 12px">
        <Page 
          :total="originData.length"         
          @on-change="changePage" 
          show-total />
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: '',
  data() {
    return {
      xyAxis: [],
      params: [],

      xValue: [],
      xOptions: [],
      yValue: null,
      yOptions: [],
      aggregateValue: 'none',
      aggregateOption: [
        {label: 'Original', value: 'none'},
        {label: 'P95', value: 'p95'},
        {label: 'Average', value: 'avg'},
        {label: 'Max', value: 'max'},
        {label: 'Min', value: 'min'}
      ],
      
      getDataParams: null,
      columns: [],
      originData: [],
      data: [],
      selectedData: []
    }
  },
  activated () {
    this.xyAxis = this.$parent.xyAxis
    this.params = this.$parent.params
    this.xOptions = JSON.parse(JSON.stringify(this.xyAxis))
    this.yOptions = JSON.parse(JSON.stringify(this.xyAxis))
  },
  methods: {
    getData () {
      this.params.forEach(p => {
        p.agg = this.aggregateValue
      })
      let params = {
        config: this.params,
        legend_x: this.xValue,
        legend_y: this.yValue
      }
      this.columns = []
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.getRData, params, (responseData) => {
        this.getDataParams = params
        this.columns.push({
          type: 'selection',
          width: 60,
          align: 'center'
        })
        responseData.data.title.forEach(i => {
          this.columns.push({
            title: i,
            key: i
          })
        })
        this.getDataParams.remove_list = []
        this.$parent.formulaParams = this.getDataParams
        this.originData = responseData.data.data
        this.data = this.originData.slice(0, 10)
      })
    },
    changePage (current) {
      this.data = this.originData.slice((current -1) * 10, current * 10)
      this.data.forEach(row => {
        row._checked = false
        if (this.selectedData.includes(row.timestamp)) {
          row._checked = true
        }
      })
    },
    selectData (selection, row) {
      this.selectedData.push(row.timestamp)
      this.getDataParams.remove_list = this.selectedData
      this.$parent.formulaParams = this.getDataParams
    },
    selectCancle (selection, row) {
      const index = this.selectedData.findIndex(i => {
        return i === row.timestamp
      })
      this.selectedData.splice(index, 1)
      this.getDataParams.remove_list = this.selectedData
      this.$parent.formulaParams = this.getDataParams
    },
    selectAllData (selection) {
      let arr = []
      selection.forEach(item => {
        arr.push(item.timestamp)
      })
      let aa = new Set(this.selectedData)
      let bb = new Set(arr)
      this.selectedData = Array.from(new Set([...aa, ...bb]))
      this.getDataParams.remove_list = this.selectedData
    },
    selectAllCancle () {
      this.data.forEach(item => {
        const index = this.selectedData.findIndex(i => {
          return i === item.timestamp
        })
        this.selectedData.splice(index, 1)
      })
    },
    clearYValue () {
      this.yValue = ''
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
.margin-left {
  margin-left: 20px;
}
</style>
