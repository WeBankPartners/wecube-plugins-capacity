<template>
  <div class="modeling-container">
    <div class="modeling-steps">
      <Row style="text-align: center">
        <Col span="8">
          <Icon type="ios-cog" class="step-icon"/>
          <span class="step-title">指标配置</span>
        </Col>
        <Col span="8">
          <Icon type="ios-construct-outline" class="step-icon"/>
          <span class="step-title">生成公式</span>
        </Col>
        <Col span="8">
          <Icon type="md-stats" class="step-icon"/>
          <span class="step-title">结果展示</span>
        </Col>
      </Row>
    </div>
    <div class="operation">
      <div class="operation-zone">
        <template v-if="current === 0">
          <ParameterConfiguration></ParameterConfiguration>
          <div class="step-control">
            <!-- <button type="button" class="btn btn-cancle-f">上一步</button> -->
            <button type="button" class="btn btn-confirm-f" @click="current = 1">下一步</button>
          </div>
        </template>
        <template v-if="current === 1">
          <DataClean :params="params" :xyAxis="xyAxis"></DataClean>
          <div class="step-control">
            <!-- <button type="button" class="btn btn-cancle-f">上一步</button> -->
            <button type="button" class="btn btn-confirm-f" @click="current = 2">生成公式</button>
          </div>
        </template>
        <template v-if="current === 2">
          <DisplayResult :formulaParams="formulaParams"></DisplayResult>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import ParameterConfiguration from '@/views/modeling/parameter-configuration'
import DataClean from '@/views/modeling/data-clean'
import DisplayResult from '@/views/modeling/display-result'
export default {
  name: '',
  data() {
    return {
      current: 0,

      params: null,
      xyAxis: null,

      formulaParams: null,
      removeList: []
    }
  },
  methods: {
    generateFormula () {
      this.current++
    },
    showResult () {
      this.current++
    },
    next () {
      if (this.current == 2) {
        this.current = 0
      } else {
        this.current += 1
      }
    }
  },
  components: {
    ParameterConfiguration,
    DataClean,
    DisplayResult
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
  height: calc(~"100vh - 140px");
  margin: 0 auto;
  background: #ffffff;
}
.step-control {
  position: fixed;
  bottom: 0;
  height: 100px;
  width: 1100px;
  padding: 34px 68px;
  border-top: 1px solid #dbe3e4;
  box-shadow: 0 -4px 4px -2px #e4e9f0;
}
button:last-child {
  float: right;
}
</style>
