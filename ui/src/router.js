import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

const router = new Router({
  scrollBehavior: () => ({ // 滚动条滚动的行为，不加这个默认就会记忆原来滚动条的位置
    y: 0
  }),
  routes: [
    {
      path: "/",
      name: "index",
      component: () => import("@/views/index"),
      title: "测试首页",
      redirect: 'capacityModeling',
      children: [
        {
          path: "/capacityModeling",
          name: "capacityModeling",
          component: () => import("@/views/capacity-modeling"),
          title: "capacityModeling"
        },
        {
          path: "/capacityForecast",
          name: "capacityForecast",
          component: () => import("@/views/capacity-forecast"),
          title: "capacityForecast"
        },
      ]
    },
    {
      path: "/test",
      name: "test",
      component: () => import("@/views/test"),
      title: "test"
    }
  ]
});

router.beforeEach((to, from, next) => {
  next()
})
export default router;
