<template>
  <Menu mode="horizontal" theme="dark" @on-select="menuChange">
    <div class="logo">
      <img src="../assets/logo.png" />
      <span>{{$t('menu.systemName')}}</span>
    </div>
    <MenuItem name="capacityModeling">
      <i class="fa fa-gears" aria-hidden="true"></i>
      {{$t("menu.capacityModeling")}}
    </MenuItem>
    <MenuItem name="capacityForecast">
      <i class="fa fa-coffee" aria-hidden="true"></i>
      {{$t("menu.capacityForecast")}}
    </MenuItem>
  </Menu>
</template>
<script>
import '@/assets/theme/dark/styls.less';
import '@/assets/theme/default/styls.less';
export default {
  data() {
    return {
      activeName: "",
      activeLang: "",
      langConfig: {
        "zh-CN": "中文",
        "en-US": "English"
      },
      lang: [{ label: "中文", value: "zh-CN" }, { label: "English", value: "en-US" }],
      username: localStorage.getItem('username')
    };
  },
  mounted() {
    if (this.langConfig[localStorage.getItem("lang")] === undefined) {
      this.activeLang = this.langConfig[
        navigator.language || navigator.userLanguage
      ];
      this.setLocale(
        navigator.language || navigator.userLanguage === "zh-CN"
          ? "zh-CN"
          : "en-US"
      );
    } else {
      this.activeLang = this.langConfig[localStorage.getItem("lang")];
    }
  },
  methods: {
    changeLang(name) {
      this.activeLang = name;
      let lang = name === "English" ? "en-US" : "zh-CN";
      this.setLocale(lang);
    },
    setLocale(lang) {
      localStorage.setItem("lang", lang);
      this.$i18n.locale = lang;
      this.$validator.locale = lang;
    },
    menuChange(name) {
      this.activeName = name;
      if (this.$route.name === name) return;
      this.$router.push({ name: name });
    }
  }
};
</script>
<style lang="less" scoped>
.ivu-menu-dark {
  background: #454a52;
}
.logo {
  float: left;
  height: inherit;
  padding-left: 30px;
  cursor: pointer;
  span {
    color: white;
    font-size: 16px;
    font-weight: bolder;
    vertical-align: top;
  }
  img {
    width: 40px;
    margin: 0 10px;
  }
}
.set-theme {
  float: right;
  width: 22px;
  height: 22px;
  border-radius: 4px;
  cursor: pointer;
  i {
    font-size: 18px;
  }
}
.menu-right {
  float: right;
  margin-right: 30px;
}
</style>
