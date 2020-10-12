<template>
  <div class="data-clean">
    <template v-if="isImportData">
      <Alert>{{$t('uploadFileTip')}}</Alert>
      <Upload 
        :on-success="importSuccess"
        :format="fileType"
        :max-size=10240
        action="/capacity/api/v1/r/excel">
        <Button icon="ios-cloud-upload-outline">Upload files</Button>
      </Upload>
    </template>
    <template>
      <Row>
        <Col span="3">
          <span class="param-title">{{$t('axisCoordinates')}}</span>
        </Col>
        <Col span="21">
          <Form :label-width="20">
            <FormItem class="param-inline">
              <Select v-model="getDataParams.legend_x" style="width:320px" :placeholder="$t('placeholder.xAxis')" filterable multiple>
                <Option v-for="item in xOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
              </Select>
            </FormItem>
            <FormItem class="param-inline">
              <Select 
                v-model="getDataParams.legend_y" 
                style="width:320px" 
                :placeholder="$t('placeholder.yAxis')" 
                filterable
                clearable
                @on-clear="clearYValue">
                <Option v-for="item in yOptions" :value="item.value" :key="item.value">{{ item.value }}</Option>
              </Select>
            </FormItem>
            <FormItem class="param-inline" v-if="!isImportData">
              <Select v-model="aggregateValue" style="width:90px" :placeholder="$t('placeholder.aggregate')" filterable>
                <Option v-for="item in aggregateOption" :value="item.value" :key="item.value">{{ item.label }}</Option>
              </Select>
            </FormItem>
          </Form>
        </Col>
      </Row>
      <Row v-if="!isImportData">
        <Col span="21" offset="3">
          <button :disabled="!(getDataParams.legend_x.length && getDataParams.legend_y)" @click="getData" type="button" class="btn btn-confirm-f margin-left">{{$t('searchData')}}</button>
        </Col>
      </Row>
    </template>
    <div style="height:500px;margin-top:40px;margin-bottom:125px">
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
      fileType: ['xls', 'xlsx'],
      xyAxis: [],
      params: [],

      xOptions: [],
      yOptions: [],
      aggregateValue: 'none',
      aggregateOption: [
        {label: 'Original', value: 'none'},
        {label: 'P95', value: 'p95'},
        {label: 'Average', value: 'avg'},
        {label: 'Max', value: 'max'},
        {label: 'Min', value: 'min'}
      ],
      
      getDataParams: {
        config: [],
        legend_x: [],
        legend_y: [],
        remove_list: []
      },
      columns: [],
      originData: [],
      data: [],
      selectedData: [],

      isImportData: false
    }
  },
  activated () {
    this.$parent.cachedCom.push(this.$vnode)
    this.isImportData = this.$parent.isImportData
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
      this.getDataParams.config = this.params
      this.columns = []
      this.$root.$capacityRequestEntrance.capacityRequestEntrance('POST', this.$root.capacityApiCenter.getRData, this.getDataParams, (responseData) => {
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
        if (this.selectedData.includes(row.id)) {
          row._checked = true
        }
      })
    },
    selectData (selection, row) {
      this.selectedData.push(Number(row.id))
      this.getDataParams.remove_list = this.selectedData
      this.$parent.formulaParams = this.getDataParams
    },
    selectCancle (selection, row) {
      const index = this.selectedData.findIndex(i => {
        return i === row.id
      })
      this.selectedData.splice(index, 1)
      this.getDataParams.remove_list = this.selectedData
      this.$parent.formulaParams = this.getDataParams
    },
    selectAllData (selection) {
      let arr = []
      selection.forEach(item => {
        arr.push(item.id)
      })
      let aa = new Set(this.selectedData)
      let bb = new Set(arr)
      this.selectedData = Array.from(new Set([...aa, ...bb]))
      this.getDataParams.remove_list = this.selectedData
    },
    selectAllCancle () {
      this.data.forEach(item => {
        const index = this.selectedData.findIndex(i => {
          return i === item.id
        })
        this.selectedData.splice(index, 1)
      })
    },
    clearYValue () {
      this.yValue = ''
    },
    importSuccess (response) {
      this.$Message.success(response.message)
      this.columns = []
      this.xyAxis = []
      this.columns.push({
        type: 'selection',
        width: 60,
        align: 'center'
      })
      response.data.table.title.forEach(val => {
        this.columns.push({
          title: val,
          key: val
        })
        if (val !== 'index') {
          this.xyAxis.push({
            label: val,
            value: val
        })
        }
      })
      this.getDataParams = {
        guid: response.data.guid,
        config: {},
        legend_x: [],
        legend_y: [],
        remove_list: []
      }
      this.$parent.formulaParams = this.getDataParams
      this.$parent.xyAxis = this.xyAxis
      this.xOptions = JSON.parse(JSON.stringify(this.xyAxis))
      this.yOptions = JSON.parse(JSON.stringify(this.xyAxis))
      this.originData = response.data.table.data
      this.data = this.originData.slice(0, 10)
    }
  },
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
