<template>
  <div class="modeling-container">
    <div style="background:white">
      <div class="modeling-steps">
        <Row>
          <Steps :current="current">
            <Step>
              <span slot="title" class="step-title">
                {{$t(steps[0])}}
              </span>
            </Step>
            <Step>
              <span slot="title" class="step-title">
                {{$t(steps[1])}}
              </span>
            </Step>
            <Step>
              <span slot="title" class="step-title">
                {{$t(steps[2])}}
              </span>
            </Step>
          </Steps>
        </Row>
      </div>
    </div>
    <div class="operation">
      <div class="operation-zone">
        <keep-alive :include="whiteList" :exclude="blackList">
          <component :is="currentComponent"></component>
        </keep-alive>
        
        <div class="step-control">
          <button v-if="current != 0" @click="upStep"  type="button" class="btn btn-cancle-f">{{$t('previous')}}:{{$t(steps[current-1])}}</button>
          <button v-if="current != 2" @click="downStep" type="button" class="btn btn-confirm-f">{{$t('nextStep')}}:{{$t(steps[current+1])}}</button>
        </div>
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
      whiteList: ['ParameterConfiguration', 'DataClean', 'DisplayResult'],
      blackList: [],
      currentComponent: 'ParameterConfiguration',
      steps:['indicatorConfiguration', 'dataCleaning', 'resultDisplay'],

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
    upStep () {
      this.current--
      this.currentComponent = this.whiteList[this.current]
    },
    downStep () {
      this.current++
      this.currentComponent = this.whiteList[this.current]
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
  width: 1100px;
  margin: 0 auto;
  padding: 16px 0;
}
.step-title {
  padding: 8px;
  font-size: 18px;
  font-weight: 500;
}
.operation {
  background: #f0f2f5;
  padding-top: 20px;
}
.operation-zone {
  width: 1100px;
  min-height: calc(~"100vh - 145px");
  margin: 0 auto;
  background: #ffffff;
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
button:last-child {
  float: right;
}
</style>
