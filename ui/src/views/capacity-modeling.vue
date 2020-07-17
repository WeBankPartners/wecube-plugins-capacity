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
          <!-- <button class="btn btn-confirm-f" @click="generateFormula">下一步</button> -->
        </template>
        <template v-if="current === 1">
          222222
          <button class="btn btn-confirm-f" @click="showResult">下一步</button>
        </template>
        <template v-if="current === 2">
          333333
        </template>
      </div>
    </div>
    
  </div>
</template>

<script>
import ParameterConfiguration from '@/views/modeling/parameter-configuration'
export default {
  name: '',
  data() {
    return {
      current: 0
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
    ParameterConfiguration
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
</style>
