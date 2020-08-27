
import capacityModeling from "@/views/capacity-modeling";
import capacityForecast from "@/views/capacity-forecast";

const router = [
  {
    path: "/capacityModeling",
    name: "capacityModeling",
    title: "容量建模",
    meta: {},
    component: capacityModeling
  },
  {
    path: "/capacityForecast",
    name: "capacityForecast",
    title: "容量预测",
    meta: {},
    component: capacityForecast
  }
];

export default router;
