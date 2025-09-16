<template>
  <div>
    <!-- 顶部统计卡片 -->
    <div class="zl">
      <div
        class="yuansu node"
        :class="nodeCollapsed ? 'collapsed' : 'expanded'"
        @click="toggleNodeCollapse"
      >
        <div class="left" :class="{ 'collapsed-left': nodeCollapsed }">
          <div class="icon-wrapper"><i class="fa fa-server icon node-icon"></i></div>
          <div class="text">Node</div>
          <div class="number">{{ stats.node }}</div>
        </div>
        <div class="right node-color"></div>
      </div>

      <div class="yuansu pod">
        <div class="left">
          <div class="icon-wrapper"><i class="fa fa-cube icon pod-icon"></i></div>
          <div class="text">Pod</div>
          <div class="number">{{ stats.pods }}</div>
        </div>
        <div class="right pod-color"></div>
      </div>

      <div class="yuansu deployment">
        <div class="left">
          <div class="icon-wrapper"><i class="fa fa-layer-group icon deployment-icon"></i></div>
          <div class="text">Deployment</div>
          <div class="number">{{ stats.deps }}</div>
        </div>
        <div class="right deployment-color"></div>
      </div>

      <div class="yuansu service">
        <div class="left">
          <div class="icon-wrapper"><i class="fa fa-network-wired icon service-icon"></i></div>
          <div class="text">Service</div>
          <div class="number">{{ stats.svcs }}</div>
        </div>
        <div class="right service-color"></div>
      </div>

      <div class="yuansu ingress">
        <div class="left">
          <div class="icon-wrapper"><i class="fa fa-route icon ingress-icon"></i></div>
          <div class="text">Ingress</div>
          <div class="number">{{ stats.ings }}</div>
        </div>
        <div class="right ingress-color"></div>
      </div>
    </div>

    <!-- 节点统计表 -->
    <div class="tb">
      <div class="tb-row" v-for="(node, index) in animatedNodes" :key="node.name">
        <div class="node-name">{{ node.name }}</div>
        <div class="node-data">
          <!-- CPU 饼图 -->
          <div class="chart-box">
            <div class="chart-inner" :ref="el => cpuChartsRefs[index] = el"></div>
            <div class="percent-label">CPU {{ node.cpuPercent.toFixed(1) }}%</div>
          </div>

          <!-- 内存 饼图 -->
          <div class="chart-box">
            <div class="chart-inner" :ref="el => memChartsRefs[index] = el"></div>
            <div class="percent-label">Mem {{ node.memPercent.toFixed(1) }}%</div>
          </div>

          <!-- CPU 使用量 -->
          <div class="data-box">
            <span>CPU用量</span>
            <b>{{ (node.cpuUsage / 1000000).toFixed(1) }} mCPU</b>
          </div>

          <!-- 内存 使用量 -->
          <div class="data-box">
            <span>内存用量</span>
            <b>{{ (node.memUsage / 1024).toFixed(1) }} Mi</b>
          </div>
        </div>
      </div>
    </div>

    <!-- 柱状图 -->
    <div class="Graph">
      <div ref="barChartContainer" class="Chart-Bar"></div>
    </div>
  </div>
</template>

<script>
import { defineComponent, reactive, ref, onMounted, nextTick, watch, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
// import axios from 'axios'
import axios from "@/main" 

export default defineComponent({
  name: 'AppDashboard',
  props: { collapsed: { type: Boolean, default: false } },
  setup(props) {
    const stats = reactive({ node: 0, pods: 0, deps: 0, svcs: 0, ings: 0 })
    const animatedNodes = reactive([])

    const cpuChartsRefs = ref([]), memChartsRefs = ref([]), barChartContainer = ref(null)
    let cpuCharts = [], memCharts = [], barChart = null
    let refreshTimer = null

    const nodeCollapsed = ref(false)
    const toggleNodeCollapse = () => { nodeCollapsed.value = !nodeCollapsed.value }

    const createRingChartOption = (used, free, usedColor, freeColor) => ({
      color: [usedColor, freeColor],
      tooltip: { trigger: 'item', formatter: '{b}: {c}%', confine: true },
      animationDuration: 800,
      series: [{
        type: 'pie',
        radius: ['70%', '90%'],
        label: { show: false },
        data: [
          { value: used, name: '已用' },
          { value: free, name: '剩余' }
        ]
      }]
    })

    const initCharts = () => {
      cpuChartsRefs.value.forEach((cpuEl, i) => {
        const memEl = memChartsRefs.value[i];
        if (!cpuEl || !memEl) return
        const cpuChart = echarts.getInstanceByDom(cpuEl) || echarts.init(cpuEl)
        const memChart = echarts.getInstanceByDom(memEl) || echarts.init(memEl)
        cpuChart.setOption(createRingChartOption(animatedNodes[i].cpuPercent, 100 - animatedNodes[i].cpuPercent, '#5b8ff9', '#e0e0e0'))
        memChart.setOption(createRingChartOption(animatedNodes[i].memPercent, 100 - animatedNodes[i].memPercent, '#61dDAA', '#e0e0e0'))
        cpuCharts.push(cpuChart); memCharts.push(memChart)
      })
    }

    const fetchStats = async () => {
      try {
        const res = await axios.get('http://192.168.216.50:8090/listall')
        const data = res.data
        stats.node = data.node || 0
        stats.pods = data.pods || 0
        stats.deps = data.deps || 0
        stats.svcs = data.svcs || 0
        stats.ings = data.ings || 0
      } catch (e) {
        console.error('获取统计数据失败:', e)
      }
    }

    const fetchNodeMetrics = async () => {
      try {
        const res = await axios.get('http://192.168.216.50:8090/nodemetrics')
        const data = res.data.data || []
        animatedNodes.splice(0, animatedNodes.length, ...data.map(n => ({
          name: n.name,
          cpuUsage: parseFloat(n.cpuUsage),         // n单位
          cpuPercent: parseFloat(n.cpuPercent),     // 百分比
          memUsage: parseFloat(n.memUsage),         // Ki单位
          memPercent: parseFloat(n.memPercent)      // 百分比
        })))
      } catch (e) {
        console.error('获取节点监控数据失败:', e)
      }
    }

    const fetchBarChartData = async () => {
      try {
        const res = await axios.get('http://192.168.216.50:8090/getnodepodstatus')
        const data = res.data.data || []
        if (!barChartContainer.value) return
        barChart = echarts.getInstanceByDom(barChartContainer.value) || echarts.init(barChartContainer.value)
        barChart.setOption({
          title: { text: 'Pod 状态分布', left: 'center' },
          tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
          legend: { top: '8%', data: ['正常', '非正常'] },
          grid: { left: '3%', right: '3%', bottom: '5%', containLabel: true },
          xAxis: { type: 'category', data: data.map(n => n.nodeName), axisLabel: { rotate: 20 } },
          yAxis: { type: 'value', name: 'Pod 数量' },
          series: [
            { name: '正常', type: 'bar', data: data.map(n => n.runningPods), itemStyle: { color: '#61dDAA' } },
            { name: '非正常', type: 'bar', data: data.map(n => n.nonRunningPods), itemStyle: { color: '#f66d6d' } }
          ]
        })
      } catch (e) {
        console.error('获取节点Pod状态失败:', e)
      }
    }

    const refreshData = async () => {
      await fetchStats()
      await fetchNodeMetrics()
      await fetchBarChartData()
      cpuCharts.forEach((c, i) => {
        if (c && animatedNodes[i]) {
          c.setOption(createRingChartOption(animatedNodes[i].cpuPercent, 100 - animatedNodes[i].cpuPercent, '#5b8ff9', '#e0e0e0'))
        }
      })
      memCharts.forEach((m, i) => {
        if (m && animatedNodes[i]) {
          m.setOption(createRingChartOption(animatedNodes[i].memPercent, 100 - animatedNodes[i].memPercent, '#61dDAA', '#e0e0e0'))
        }
      })
    }

    onMounted(async () => {
      await fetchStats()
      await fetchNodeMetrics()
      await nextTick(); 
      initCharts()
      await fetchBarChartData()
      refreshTimer = setInterval(refreshData, 3000)
    })

    onBeforeUnmount(() => {
      cpuCharts.forEach(c => c.dispose())
      memCharts.forEach(m => m.dispose())
      if (barChart) barChart.dispose()
      cpuCharts = []
      memCharts = []
      barChart = null
      if (refreshTimer) clearInterval(refreshTimer)
    })

    watch(() => props.collapsed, () => {
      nextTick(() => {
        setTimeout(() => {
          cpuCharts.forEach(c => c.resize())
          memCharts.forEach(m => m.resize())
          if (barChart) requestAnimationFrame(() => barChart.resize())
        }, 350)
      })
    })

    return { stats, animatedNodes, cpuChartsRefs, memChartsRefs, barChartContainer, nodeCollapsed, toggleNodeCollapse }
  }
})
</script>

<style scoped>
/* 保持原有样式 */
.zl{display:flex;justify-content:space-between;align-items:center;height:120px;padding:0 20px;background:#f9f9f9;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,0.05);}
.yuansu{width:140px;height:100%;display:flex;justify-content:space-between;align-items:center;padding:10px;border-radius:8px;box-shadow:0 1px 4px rgba(0,0,0,0.1);cursor:pointer;flex-shrink:0;transition:width 0.3s ease;}
.yuansu.collapsed{width:40px;}
.left{display:flex;flex-direction:column;gap:5px;color:#333;min-width:24px;overflow:hidden;transition:all 0.3s ease;}
.left.collapsed-left .text,.left.collapsed-left .number{opacity:0;width:0;height:0;}
.icon-wrapper{width:24px;height:24px;display:flex;align-items:center;justify-content:center;overflow:hidden;}
.icon{font-size:24px;transition:transform 0.3s ease;}
.node-icon{color:#5b8ff9;}
.pod-icon{color:#61dDAA;}
.deployment-icon{color:#65789b;}
.service-icon{color:#f6c022;}
.ingress-icon{color:#f66d6d;}
.yuansu:hover .icon{transform:rotate(15deg) scale(1.2);}
.right{width:12px;height:100%;border-radius:8px;}
.node-color{background-color:#5b8ff9;}
.pod-color{background-color:#61dDAA;}
.deployment-color{background-color:#65789b;}
.service-color{background-color:#f6c022;}
.ingress-color{background-color:#f66d6d;}

.tb{padding:20px;background:#f9f9f9;border-radius:8px;overflow-x:auto;}
.tb-row{display:flex;align-items:center;margin-bottom:15px;padding:10px 15px;background:#fff;border-radius:8px;box-shadow:0 1px 4px rgba(0,0,0,0.08);flex-wrap:wrap;transition:transform 0.2s,box-shadow 0.2s;}
.tb-row:hover{transform:translateY(-3px);box-shadow:0 4px 12px rgba(0,0,0,0.12);}
.node-name{width:180px;font-weight:bold;color:#333;}
.node-data{display:flex;flex-wrap:wrap;gap:15px;align-items:center;flex:1;}
.chart-box{width:80px;height:80px;position:relative;flex-shrink:0;min-width:60px;min-height:60px;z-index:10;}
.chart-inner{width:100%;height:100%;}
.percent-label{position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);font-size:12px;font-weight:bold;color:#333;pointer-events:none;}
.data-box{padding:10px 15px;background:#f9f9f9;border-radius:6px;box-shadow:0 1px 3px rgba(0,0,0,0.08);transition:transform 0.2s,box-shadow 0.2s;}
.data-box:hover{transform:translateY(-2px);box-shadow:0 3px 8px rgba(0,0,0,0.12);}
.data-box span{font-size:12px;color:#666;display:block;margin-bottom:5px;}
.data-box b{font-size:14px;color:#333;}

.Graph{width:100%;height:300px;margin-top:20px;padding:10px;border-radius:8px;background:#fff;box-shadow:0 1px 6px rgba(0,0,0,0.1);}
.Chart-Bar {width:100%;height:300px;min-height:200px;}
:deep(.echarts-tooltip) {z-index: 9999 !important;}
</style>
