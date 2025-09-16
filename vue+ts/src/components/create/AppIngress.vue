<template>
  <div class="ingress-create-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-plus-circle icon"></i> 创建 Ingress
      </h2>

      <!-- 只读的基本固定字段 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> networking.k8s.io/v1</p>
        <p><strong>kind：</strong> Ingress</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- Ingress 基本信息 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-info-circle section-icon blue"></i> 基本信息
            </h3>
            <div class="form-item">
              <label>名称：</label>
              <input v-model="ingress.name" type="text" placeholder="请输入 Ingress 名称" />
            </div>
            <div class="form-item">
              <label>命名空间：</label>
              <select v-model="ingress.namespace">
                <option value="">请选择命名空间</option>
                <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
              </select>
            </div>
           <div class="form-item">
  <label>IngressClass：</label>
  <input
    list="ingress-class-list"
    v-model="ingress.ingressClassName"
    placeholder="请选择或输入 IngressClass"
  />
  <datalist id="ingress-class-list">
    <option v-for="cls in ingressClasses" :key="cls" :value="cls">{{ cls }}</option>
  </datalist>
</div>
          </div>

          <!-- Annotations -->
          <div class="form-section">
            <h3>
              <i class="fa fa-sticky-note section-icon orange"></i> Annotations
            </h3>
            <div v-for="(annotation, idx) in ingress.annotations" :key="'anno-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="annotation.key" type="text" placeholder="如 nginx.ingress.kubernetes.io/rewrite-target" />
              <label>Value：</label>
              <input v-model="annotation.value" type="text" placeholder="如 /" />
              <button class="btn btn-danger btn-small" @click="removeAnnotation(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addAnnotation">
              <i class="fa fa-plus"></i> 添加 Annotation
            </button>
          </div>

          <!-- Rules 配置 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-link section-icon green"></i> Rules 配置
            </h3>
            <div v-for="(rule, idx) in ingress.rules" :key="'rule-' + idx" class="rule-item">
              <div class="form-item">
                <label>Host：</label>
                <input v-model="rule.host" type="text" placeholder="如 nginx.example.com" />
              </div>
              <div v-for="(pathObj, pIdx) in rule.paths" :key="'path-' + pIdx" class="path-item">
                <div class="path-row">
                  <label>路径：</label>
                  <input v-model="pathObj.path" type="text" placeholder="如 /" />
                </div>
                <div class="path-row">
                  <label>PathType：</label>
                  <select v-model="pathObj.pathType">
                    <option value="Prefix">Prefix</option>
                    <option value="Exact">Exact</option>
                    <option value="ImplementationSpecific">ImplementationSpecific</option>
                  </select>
                </div>
                <div class="path-row">
                  <label>服务：</label>
                  <input v-model="pathObj.serviceName" type="text" placeholder="如 nginx-service" />
                </div>
                <div class="path-row">
                  <label>端口：</label>
                  <input v-model="pathObj.servicePort" type="number" placeholder="如 80" />
                </div>
                <button class="btn btn-danger btn-small" @click="removePath(idx, pIdx)">
                  <i class="fa fa-trash"></i> 移除路径
                </button>
              </div>
              <button class="btn btn-add btn-small" @click="addPath(idx)">
                <i class="fa fa-plus"></i> 添加路径
              </button>
              <button class="btn btn-danger btn-small" style="margin-top:6px;" @click="removeRule(idx)">
                <i class="fa fa-trash"></i> 移除 Host
              </button>
            </div>
            <button class="btn btn-add" @click="addRule">
              <i class="fa fa-plus"></i> 添加 Host
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
        <button class="btn btn-primary" @click="createIngress">
          <i class="fa fa-check"></i> 创建 Ingress
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
  name: "AppIngressCreate",
  setup() {
    const { namespaces } = useK8sData()
    const ingressClasses = ref(["nginx", "traefik", "istio"])

    const ingress = ref({
      name: "",
      namespace: "",
      ingressClassName: "",
      annotations: [{ key: "", value: "" }],
      rules: [
        {
          host: "",
          paths: [{ path: "", pathType: "Prefix", serviceName: "", servicePort: null }]
        }
      ]
    })

    const addAnnotation = () => ingress.value.annotations.push({ key: "", value: "" })
    const removeAnnotation = idx => ingress.value.annotations.splice(idx, 1)

    const addRule = () => ingress.value.rules.push({
      host: "",
      paths: [{ path: "", pathType: "Prefix", serviceName: "", servicePort: null }]
    })
    const removeRule = idx => ingress.value.rules.splice(idx, 1)

    const addPath = (ruleIdx) => ingress.value.rules[ruleIdx].paths.push({
      path: "", pathType: "Prefix", serviceName: "", servicePort: null
    })
    const removePath = (ruleIdx, pathIdx) => ingress.value.rules[ruleIdx].paths.splice(pathIdx, 1)

    const resetForm = () => {
      ingress.value = {
        name: "",
        namespace: "",
        ingressClassName: "",
        annotations: [{ key: "", value: "" }],
        rules: [{ host: "", paths: [{ path: "", pathType: "Prefix", serviceName: "", servicePort: null }] }]
      }
    }

    const yamlOutput = computed(() => {
      const obj = {
        apiVersion: "networking.k8s.io/v1",
        kind: "Ingress",
        metadata: {
          name: ingress.value.name || "",
          namespace: ingress.value.namespace || "",
          annotations: ingress.value.annotations.reduce((acc, cur) => {
            if (cur.key) acc[cur.key] = cur.value
            return acc
          }, {})
        },
        spec: {
          ingressClassName: ingress.value.ingressClassName || "",
          rules: ingress.value.rules.map(r => ({
            host: r.host || "",
            http: {
              paths: r.paths.map(p => ({
                path: p.path || "",
                pathType: p.pathType || "Prefix",
                backend: {
                  service: {
                    name: p.serviceName || "",
                    port: { number: Number(p.servicePort) || 80 }
                  }
                }
              }))
            }
          }))
        }
      }
      return yaml.dump(obj, { noRefs: true })
    })

   const createIngress = async () => {
  if (!ingress.value.name || !ingress.value.namespace) {
    alert("请填写名称和命名空间")
    return
  }

  // 构造请求 payload，符合后端接口 ReqIngressData
  const payload = {
    name: ingress.value.name,
    namespace: ingress.value.namespace,
    ingressClassName: ingress.value.ingressClassName,
    annotations: ingress.value.annotations.reduce((acc, cur) => {
      if (cur.key) acc[cur.key] = cur.value
      return acc
    }, {}),
    rules: ingress.value.rules.map(r => ({
      host: r.host,
      paths: r.paths.map(p => ({
        path: p.path,
        pathType: p.pathType,
        serviceName: p.serviceName,
        servicePort: p.servicePort
      }))
    }))
  }

  try {
    const res = await axios.post(
      "http://192.168.216.50:8090/ingest/CreateIngress",
      payload
    )
    if (res.data && res.data.message) {
      alert("创建结果: " + res.data.message)
    } else {
      alert("创建成功")
    }
  } catch (err) {
    console.error("创建 Ingress 失败:", err)
    alert("创建失败，请检查控制台信息")
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
      ingressClasses,
      ingress,
      addAnnotation,
      removeAnnotation,
      addRule,
      removeRule,
      addPath,
      removePath,
      resetForm,
      createIngress,
      yamlOutput,
      copyYaml
    }
  }
})
</script>

<style scoped>
.title>.icon { color: #9c27b0 }
.page-container { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.readonly-info { background: #f5f5f5; padding: 8px 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px; }
.main-content { display: flex; gap: 20px; }
.form-content { flex: 1.2 }
.yaml-preview { flex: 1; background: #fff; color: #333; padding: 16px; border-radius: 8px; font-size: 15px; overflow-x: auto; border: 2px solid transparent; background-clip: padding-box; position: relative; }
.yaml-preview::before { content: ""; position: absolute; inset: 0; border-radius: 8px; padding: 2px; background: linear-gradient(45deg, #42a5f5, #66bb6a); -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0); -webkit-mask-composite: xor; mask-composite: exclude; }
.copy-btn { position: absolute; top: 8px; right: 12px; background: #42a5f5; color: #fff; border: none; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.form-section { margin-bottom: 20px; padding: 15px; border: 1px solid #eee; border-radius: 10px; }
.section-icon.blue { color: #42a5f5 }
.section-icon.green { color: #66bb6a }
.section-icon.orange { color: #ff9800 }
.form-item { display: flex; align-items: center; margin: 10px 0; gap: 8px; flex-wrap: wrap; }
.form-item label { width: 120px; font-weight: 500; color: #555; }
.form-item input, .form-item select { flex: 1; padding: 6px 10px; border: 1px solid #ddd; border-radius: 6px; }
.rule-item { flex-direction: column; border: 1px dashed #ddd; padding: 10px; border-radius: 8px; margin-bottom: 10px; }
.path-item { display: flex; flex-direction: column; gap: 6px; margin-bottom: 6px; }
.path-row { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }
.form-actions { margin-top: 20px; display: flex; gap: 12px }
.btn { padding: 6px 12px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.btn-primary { background: #9c27b0; color: #fff }
.btn-secondary { background: #ccc; color: #333 }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px }
.btn-danger { background: #f66d6d; color: #fff }
.btn-small { font-size: 12px; padding: 4px 8px }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px; }
</style>
