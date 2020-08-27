# Wecube-Plugins-Capacity 容量插件

<p align="left">
    <a href="https://opensource.org/licenses/Apache-2.0" alt="License">
        <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" /></a>
    <a href="#" alt="Java">
        <img src="https://img.shields.io/badge/language-go-orange.svg" /></a>
    <a href="#" alt="Vue">
        <img src="https://img.shields.io/badge/language-vue-green.svg" /></a>
</p>

## 简介
容量插件是一套集成监控数据多元回归分析建模和数据预测功能的插件系统，它可用于分析业务依赖关系、查找业务异常区间和预测业务趋势，从而来进行对资源的合理规划与改进。  

## 主要功能
容量插件包括以下功能：

- 数据建模: 支持通过选取Open-Monitor的监控数据或Excel文件上传的方式来配置数据，并从中选取一个或多个自变量和一个因变量来做线性回归分析，生成多元一次的线性回归模型公式；
- 容量预测: 支持对建模公式增加多个自定义的自变量来生成对应的因变量；
- 模型对比: 支持对建模结果进行二次建模并分析对比；

## 快速入门
程序采用容器化部署。

如何编译，请查看以下文档
[Wecube-Plugins-Capacity编译文档](doc/compile_guide.md)

如何安装， 请查看以下文档
[Wecube-Plugins-Capacity部署文档](doc/install_guide.md)


## 版权声明
WeCMDB是基于 Apache License 2.0 协议， 详情请参考
[LICENSE](LICENSE)

## 社区
- 如果您想得到最快的响应，请给我们提[Issue](https://github.com/WeBankPartners/wecube-plugins-capacity/issues/new/choose)或扫描下面的二维码，我们会第一时间反馈。

	<div align="left">
	<img src="doc/images/wecube_qr_code.png"  height="200" width="200">
	</div>


- 联系我们：fintech@webank.com