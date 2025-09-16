<template>
  <div class="service-create-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-network-wired icon"></i> 创建 Service
      </h2>

      <!-- 只读的基本固定字段 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> v1</p>
        <p><strong>kind：</strong> Service</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- 基本信息 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-info-circle section-icon blue"></i> 基本信息
            </h3>
            <div class="form-item">
              <label>名称：</label>
              <input v-model="service.name" type="text" placeholder="请输入 Service 名称" />
            </div>
            <div class="form-item">
              <label>命名空间：</label>
              <select v-model="service.namespace">
                <option value="">请选择命名空间</option>
                <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
              </select>
            </div>
          </div>

          <!-- Labels -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Service Labels
            </h3>
            <div v-for="(label, idx) in service.labels" :key="'label-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="label.key" type="text" placeholder="如 app" />
              <label>Value：</label>
              <input v-model="label.value" type="text" placeholder="如 nginx" />
              <button class="btn btn-danger btn-small" @click="removeLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addLabel">
              <i class="fa fa-plus"></i> 添加 Label
            </button>
          </div>

          <!-- Selector -->
          <div class="form-section">
            <h3>
              <i class="fa fa-filter section-icon green"></i> Selector
            </h3>
            <p class="labels-hint">⚠️ Selector 用于匹配 Pod 标签</p>
            <div v-for="(sel, idx) in service.selectors" :key="'sel-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="sel.key" type="text" placeholder="如 app" />
              <label>Value：</label>
              <input v-model="sel.value" type="text" placeholder="如 nginx" />
              <button class="btn btn-danger btn-small" @click="removeSelector(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addSelector">
              <i class="fa fa-plus"></i> 添加 Selector
            </button>
          </div>

          <!-- 类型 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-cogs section-icon purple"></i> Service 类型
            </h3>
            <div class="form-item">
              <label>类型：</label>
              <select v-model="service.type">
                <option value="ClusterIP">ClusterIP</option>
                <option value="NodePort">NodePort</option>
                <option value="LoadBalancer">LoadBalancer</option>
              </select>
            </div>
          </div>

          <!-- 端口配置（两行展示） -->
          <div class="form-section">
            <h3>
              <i class="fa fa-plug section-icon teal"></i> 端口配置
            </h3>
            <div v-for="(port, idx) in service.ports" :key="'port-' + idx" class="port-item">
              <!-- 第一行：协议 + Service 对外端口 -->
              <div class="port-row">
                <label>协议：</label>
                <select v-model="port.protocol">
                  <option value="TCP">TCP</option>
                  <option value="UDP">UDP</option>
                </select>
                <label>端口：</label>
                <input v-model="port.port" type="number" placeholder="Service 对外端口" />
              </div>

              <!-- 第二行：TargetPort + 删除按钮 -->
              <div class="port-row">
                <label>TargetPort：</label>
                <input v-model="port.targetPort" type="number" placeholder="Pod 容器端口" />
                <button class="btn btn-danger btn-small" @click="removePort(idx)">
                  <i class="fa fa-trash"></i>
                </button>
              </div>
            </div>
            <button class="btn btn-add" @click="addPort">
              <i class="fa fa-plus"></i> 添加端口
            </button>
          </div>
        </div>

        <!-- 右边 YAML 实时预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> YAML 预览</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre>{{ yamlOutput }}</pre>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <button class="btn btn-primary" @click="createService">
          <i class="fa fa-check"></i> 创建 Service
        </button>
        <button class="btn btn-secondary" @click="resetForm">
          <i class="fa fa-refresh"></i> 重置
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed } from "vue"
// import axios from "axios"
import axios from "@/main" 
import { useK8sData } from "@/js/useK8sData"
import yaml from "js-yaml"

export default defineComponent({
  name: "AppServiceCreate",
  setup() {
    const { namespaces } = useK8sData()

    // 默认值为空
    const service = ref({
      name: "",
      namespace: "",
      labels: [{ key: "", value: "" }],
      selectors: [{ key: "", value: "" }],
      type: "ClusterIP",
      ports: [{ protocol: "TCP", port: "", targetPort: "" }]
    })

    const labelsObj = computed(() => {
      const obj = {}
      service.value.labels.forEach(l => { if (l.key && l.value) obj[l.key] = l.value })
      return obj
    })

    const selectorsObj = computed(() => {
      const obj = {}
      service.value.selectors.forEach(s => { if (s.key && s.value) obj[s.key] = s.value })
      return obj
    })

    const yamlOutput = computed(() => {
      const obj = {
        apiVersion: "v1",
        kind: "Service",
        metadata: {
          name: service.value.name,
          namespace: service.value.namespace,
          labels: labelsObj.value
        },
        spec: {
          selector: selectorsObj.value,
          type: service.value.type,
          ports: service.value.ports.map(p => ({
            protocol: p.protocol,
            port: p.port ? Number(p.port) : "",
            targetPort: p.targetPort ? Number(p.targetPort) : ""
          }))
        }
      }
      return yaml.dump(obj, { noRefs: true })
    })

    const addLabel = () => service.value.labels.push({ key: "", value: "" })
    const removeLabel = idx => service.value.labels.splice(idx, 1)
    const addSelector = () => service.value.selectors.push({ key: "", value: "" })
    const removeSelector = idx => service.value.selectors.splice(idx, 1)
    const addPort = () => service.value.ports.push({ protocol: "TCP", port: "", targetPort: "" })
    const removePort = idx => service.value.ports.splice(idx, 1)

    const resetForm = () => {
      service.value = {
        name: "",
        namespace: "",
        labels: [{ key: "", value: "" }],
        selectors: [{ key: "", value: "" }],
        type: "ClusterIP",
        ports: [{ protocol: "TCP", port: "", targetPort: "" }]
      }
    }

  const createService = async () => {
  try {
    // 调整为新的接口
    const res = await axios.post(
      "http://192.168.216.50:8090/service/CreateService",
      {
        name: service.value.name,
        namespace: service.value.namespace,
        labels: labelsObj.value,
        selectors: selectorsObj.value,
        type: service.value.type,
        ports: service.value.ports.map(p => ({
          protocol: p.protocol,
          port: Number(p.port),
          targetPort: Number(p.targetPort)
        }))
      }
    )
    if (res.data && res.data.message) {
      alert(res.data.message) // 成功或错误信息都提示
    } else {
      alert("创建 Service 成功")
    }
  } catch (err) {
    console.error(err)
    alert("创建 Service 失败")
  }
}

    const copyYaml = async () => {
      try {
        await navigator.clipboard.writeText(yamlOutput.value)
        alert("已复制到剪贴板")
      } catch {
        alert("复制失败，请手动选择复制")
      }
    }

    return {
      namespaces,
      service,
      yamlOutput,
      addLabel,
      removeLabel,
      addSelector,
      removeSelector,
      addPort,
      removePort,
      resetForm,
      createService,
      copyYaml
    }
  }
})
</script>

<style scoped>
.title>.icon { color: #ff9800 }
.page-container { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.readonly-info { background: #f5f5f5; padding: 8px 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px; }
.main-content { display: flex; gap: 20px; }
.form-content { flex: 1 }
.yaml-preview {
  flex: 1;
  background: #fff;
  color: #333;
  padding: 16px;
  border-radius: 8px;
  font-size: 15px;
  overflow-x: auto;
  border: 2px solid transparent;
  background-clip: padding-box;
  position: relative;
}
.yaml-preview::before {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: 8px;
  padding: 2px;
  background: linear-gradient(45deg, #42a5f5, #66bb6a);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
          mask-composite: exclude;
}
.copy-btn { position: absolute; top: 8px; right: 12px; background: #42a5f5; color: #fff; border: none; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.form-section { margin-bottom: 20px; padding: 15px; border: 1px solid #eee; border-radius: 10px; }
.section-icon.blue { color: #42a5f5 }
.section-icon.green { color: #66bb6a }
.section-icon.orange { color: #ff9800 }
.section-icon.purple { color: #9c27b0 }
.container-icon { color: #6a1b9a; margin-right: 6px; }
.form-item { display: flex; align-items: center; margin: 10px 0; gap: 8px; }
.form-item label { width: 60px; font-weight: 500; color: #555; }
.form-item input, .form-item select { flex: 1; padding: 6px 10px; border: 1px solid #ddd; border-radius: 6px; }
.container-card { background: #fafafa; border: 1px solid #ddd; border-radius: 10px; padding: 12px; margin-bottom: 12px; }
.form-actions { margin-top: 20px; display: flex; gap: 12px }
.btn { padding: 6px 12px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.btn-primary { background: #4CAF50; color: #fff }
.btn-secondary { background: #ccc; color: #333 }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px }
.btn-danger { background: #f66d6d; color: #fff }
.btn-small { font-size: 12px; padding: 4px 8px }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px; }

/* 端口配置两行布局 */
.port-item { display: flex; flex-direction: column; gap: 6px; }
.port-row { display: flex; align-items: center; gap: 8px; }
.port-row label { width: 80px; }
</style>
